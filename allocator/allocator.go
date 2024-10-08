package main

import (
	"context"
	pb "gitbub.com/uberstrate/idl"
	"log"
	"sort"
	"time"
)

type Allocator struct {
	client pb.StateStoreServiceClient
}

func NewAllocator(client pb.StateStoreServiceClient) *Allocator {
	return &Allocator{
		client: client,
	}
}

// Map to convert string to enum
var podPhaseMap = map[string]pb.PodPhase{
	"PENDING_NODE_ASSIGNMENT": pb.PodPhase_PENDING_NODE_ASSIGNMENT,
	"NODE_ASSIGNED":           pb.PodPhase_NODE_ASSIGNED,
	"RUNNING":                 pb.PodPhase_RUNNING,
	"FAILED":                  pb.PodPhase_FAILED,
}

func (s *Allocator) getPods(all bool, phase string) ([]*pb.Pod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := s.client.GetPods(ctx, &pb.GetPodRequest{All: all, Phase: podPhaseMap[phase]})
	if err != nil {
		return nil, err
	}
	return r.Pods, nil
}

func (s *Allocator) getNodes(minGeneration int64) ([]*pb.Node, error) {
	//log.Printf("Getting nodes with min_generation: %d", minGeneration)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := s.client.GetNodes(ctx, &pb.GetNodeRequest{AboveGenerationNumber: minGeneration})
	if err != nil {
		return nil, err
	}
	return r.Nodes, nil
}

func (s *Allocator) start() {
	uptimeTicker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-uptimeTicker.C:
			s.loop()
		}
	}
}

func (s *Allocator) loop() {
	nodes, err := s.getNodes(0)
	if err != nil {
		log.Printf("error calling function GetNodes: %v", err)
	} else {
		//log.Printf("Response from gRPC server's GetNodes total node: %d", len(nodes))
	}

	//log.Printf("Running reconcilliation loop")
	podAlls, err := s.getPods(true /* all */, "" /* phase */)
	if err != nil {
		log.Printf("error calling function GetPods: %v", err)
		return
	}
	//log.Printf("Response from gRPC server's GetPods total pod: %d", len(podAlls))

	schedulablePods := make([]*pb.Pod, 0, len(podAlls))
	for _, pod := range podAlls {
		if pod.Status.Phase == pb.PodPhase_PENDING_NODE_ASSIGNMENT ||
			pod.Status.Condition == pb.PodCondition_REALLOCATION_REQUIRED {
			schedulablePods = append(schedulablePods, pod)
		}
	}

	log.Printf("Current state:")
	for _, pod := range podAlls {
		log.Printf("	Pod:%s running on node:%s", pod.Metadata.Name, pod.Status.NodeUuid)
	}

	if len(schedulablePods) > 0 {
		log.Printf("Needs scheduling for %d pods.", len(schedulablePods))
		for _, pod := range schedulablePods {
			log.Printf("	Pod name: %s", pod.Metadata.Name)
		}
	}

	// Sort schedulablePods by priority in descending order
	sort.Slice(schedulablePods, func(i, j int) bool {
		return schedulablePods[i].Spec.Priority > schedulablePods[j].Spec.Priority
	})

	availableResourcesOnNode := make(map[string]*pb.Resource)
	for _, node := range nodes {
		availableResourcesOnNode[node.Metadata.Uuid] = availableResources(node, podAlls)
	}

	schedulePods := make([]*pb.Pod, 0, len(podAlls))
	for _, pod := range schedulablePods {
		log.Printf("Scheduling pod: %s", pod.Metadata.Name)
		feasibleNodes := make([]*pb.Node, 0, len(nodes))
		for _, node := range nodes {
			// Check if Node meets pod placement requirements
			if !passNodeSelector(node, pod) {
				//log.Printf("Node: %s does not meet node selector requirements for pod: %s", node.Metadata.Name, pod.Metadata.Name)
				continue
			}

			// Check if Node has enough resources
			availableResources := availableResourcesOnNode[node.Metadata.Uuid]
			//log.Printf("Node: %s has available resources: %+v", node.Metadata.Name, availableResources)
			if availableResources.Cpu < pod.Spec.ResourceRequirement.Cpu ||
				availableResources.Ram < pod.Spec.ResourceRequirement.Ram ||
				availableResources.Storage < pod.Spec.ResourceRequirement.Storage {
				continue
			}
			feasibleNodes = append(feasibleNodes, node)
		}
		//log.Printf("Found %d feasible nodes for pod: %s", len(feasibleNodes), pod.Metadata.Name)
		bestFitNode := findBestFitNode(feasibleNodes)
		if bestFitNode == nil {
			log.Printf("No feasible node found for pod: %s", pod.Metadata.Name)
			continue
		}
		log.Printf("Assigning pod: %s to node: %s", pod.Metadata.Name, bestFitNode.Metadata.Name)

		// Update available resources on the node
		availableResourcesOnNode[bestFitNode.Metadata.Uuid].Cpu -= pod.Spec.ResourceRequirement.Cpu
		availableResourcesOnNode[bestFitNode.Metadata.Uuid].Ram -= pod.Spec.ResourceRequirement.Ram
		availableResourcesOnNode[bestFitNode.Metadata.Uuid].Storage -= pod.Spec.ResourceRequirement.Storage
		pod.Status.NodeUuid = bestFitNode.Metadata.Uuid
		pod.Status.Phase = pb.PodPhase_NODE_ASSIGNED
		pod.Status.Condition = pb.PodCondition_UNKNOWN
		for i, v := range podAlls {
			if v == pod {
				podAlls[i] = pod
				break
			}
		}
		schedulePods = append(schedulePods, pod)
	}

	_, err = s.client.UpdatePods(context.Background(), &pb.UpdatePodRequest{Pods: schedulePods})
	if err != nil {
		log.Printf("error calling function UpdatePods: %v", err)
	}
}

func findBestFitNode(nodes []*pb.Node) *pb.Node {
	if len(nodes) == 0 {
		return nil
	}
	bestFitNode := nodes[0]
	for _, node := range nodes {
		// This can be improved by using a better algorithm.
		// For now, we are just using the node with the least CPU capacity.
		// This is not ideal for when multiple allocator are running since
		//all of them would pick same node.
		if node.Status.Capacity.Cpu < bestFitNode.Status.Capacity.Cpu {
			bestFitNode = node
		}
	}
	return bestFitNode
}

func passNodeSelector(node *pb.Node, pod *pb.Pod) bool {
	for key, valueToMatch := range pod.Spec.NodeSelectorLabels {
		if key == "" || valueToMatch == "" {
			continue
		}
		value, exists := node.Metadata.Labels[key]
		if !exists || value != valueToMatch {
			//log.Printf("Node: %s does not meet pod selector "+
			//	"requirements for pod: %s. Didn't find key: %s value: %s",
			//	node.Metadata.Name, pod.Metadata.Name, key, valueToMatch)
			return false
		}
	}
	return true
}

func availableResources(node *pb.Node, pods []*pb.Pod) *pb.Resource {
	if node.Status.Tainted != "" {
		return &pb.Resource{
			Cpu:     0,
			Ram:     0,
			Storage: 0,
		}
	}

	resources := &pb.Resource{
		Cpu:     node.Status.Capacity.Cpu,
		Ram:     node.Status.Capacity.Ram,
		Storage: node.Status.Capacity.Storage,
	}
	// Calculate the available resources on a node
	for _, pod := range pods {
		if pod.Status.NodeUuid == node.Metadata.Uuid && pod.Status.Phase == pb.PodPhase_RUNNING {
			// Deduct the resources used by the running pods
			resources.Cpu -= pod.Spec.ResourceRequirement.Cpu
			resources.Ram -= pod.Spec.ResourceRequirement.Ram
			resources.Storage -= pod.Spec.ResourceRequirement.Storage
		}
	}
	return resources
}

func (s *Allocator) Run() {
	log.Printf("Starting allocator")
	go s.start()
}
