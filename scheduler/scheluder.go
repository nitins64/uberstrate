package main

import (
	"context"
	pb "gitbub.com/uberstrate/idl"
	"log"
	"sort"
	"time"
)

type Scheduler struct {
	client pb.StateStoreServiceClient
}

func NewScheduler(client pb.StateStoreServiceClient) *Scheduler {
	return &Scheduler{
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

func (s *Scheduler) getPods(all bool, phase string) ([]*pb.Pod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := s.client.GetPods(ctx, &pb.GetPodRequest{All: all, Phase: podPhaseMap[phase]})
	if err != nil {
		return nil, err
	}
	return r.Pods, nil
}

func (s *Scheduler) getNodes(minGeneration int64) ([]*pb.Node, error) {
	log.Printf("Getting nodes with min_generation: %d", minGeneration)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := s.client.GetNodes(ctx, &pb.GetNodeRequest{AboveGenerationNumber: minGeneration})
	if err != nil {
		return nil, err
	}
	return r.Nodes, nil
}

func (s *Scheduler) start() {
	uptimeTicker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-uptimeTicker.C:
			s.loop()
		}
	}
}

func (s *Scheduler) loop() {
	nodes, err := s.getNodes(0)
	if err != nil {
		log.Printf("error calling function GetNodes: %v", err)
	} else {
		log.Printf("Response from gRPC server's GetNodes total node: %d", len(nodes))
	}

	log.Printf("Getting all podAlls")
	podAlls, err := s.getPods(true /* all */, "" /* phase */)
	if err != nil {
		log.Fatalf("error calling function GetPods: %v", err)
	}
	log.Printf("Response from gRPC server's GetPods total pod: %d", len(podAlls))

	schedulablePods := make([]*pb.Pod, 0, len(podAlls))
	for _, pod := range podAlls {
		if pod.Status.Phase == pb.PodPhase_PENDING_NODE_ASSIGNMENT {
			schedulablePods = append(schedulablePods, pod)
		}
	}

	log.Printf("Scheduling %d podAlls", len(schedulablePods))

	// Sort schedulablePods by priority in descending order
	sort.Slice(schedulablePods, func(i, j int) bool {
		return schedulablePods[i].Spec.Priority > schedulablePods[j].Spec.Priority
	})

	schedudedPods := make([]*pb.Pod, 0, len(podAlls))
	for _, pod := range schedulablePods {
		log.Printf("Scheduling pod: %s", pod.Metadata.Name)
		feasibleNodes := make([]*pb.Node, 0, len(nodes))
		for _, node := range nodes {
			// Check if Node meets pod placement requirements
			if !passNodeSelector(node, pod) {
				log.Printf("Node: %s does not meet node selector requirements for pod: %s", node.Metadata.Name, pod.Metadata.Name)
				continue
			}

			// Check if Node has enough resources
			avaiableResources := availableResources(node, podAlls)
			log.Printf("Node: %s has available resources: %+v", node.Metadata.Name, avaiableResources)
			if avaiableResources.Cpu < pod.Spec.ResourceRequirement.Cpu ||
				avaiableResources.Ram < pod.Spec.ResourceRequirement.Ram ||
				avaiableResources.Storage < pod.Spec.ResourceRequirement.Storage {
				continue
			}
			feasibleNodes = append(feasibleNodes, node)
		}
		log.Printf("Found %d feasible nodes for pod: %s", len(feasibleNodes), pod.Metadata.Name)
		bestFitNode := findBestFitNode(feasibleNodes)
		if bestFitNode == nil {
			log.Printf("No feasible node found for pod: %s", pod.Metadata.Name)
			continue
		}
		log.Printf("Assigning pod: %s to node: %s", pod.Metadata.Name, bestFitNode.Metadata.Name)
		pod.Status.NodeUuid = bestFitNode.Metadata.Uuid
		pod.Status.Phase = pb.PodPhase_NODE_ASSIGNED
		for i, v := range podAlls {
			if v == pod {
				podAlls[i] = pod
				break
			}
		}
		schedudedPods = append(schedudedPods, pod)
	}

	_, err = s.client.UpdatePods(context.Background(), &pb.UpdatePodRequest{Pods: schedudedPods})
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
		// This is not ideal for when multiple scheduler are running since
		//all of them would pick same node.
		if node.Status.Capacity.Cpu < bestFitNode.Status.Capacity.Cpu {
			bestFitNode = node
		}
	}
	return bestFitNode
}

func passNodeSelector(node *pb.Node, pod *pb.Pod) bool {
	for key, valueToMatch := range pod.Spec.NodeSelectorLabels {
		value, exists := node.Metadata.Labels[key]
		if !exists || value != valueToMatch {
			log.Printf("Node: %s does not meet pod selector "+
				"requirements for pod: %s. Didn't find key: %s value: %s",
				node.Metadata.Name, pod.Metadata.Name, key, valueToMatch)
			return false
		}
	}
	return true
}

func availableResources(node *pb.Node, pods []*pb.Pod) *pb.Resource {
	resources := &pb.Resource{
		Cpu:     node.Status.Capacity.Cpu,
		Ram:     node.Status.Capacity.Ram,
		Storage: node.Status.Capacity.Storage,
	}
	// Calculate the available resources on a node
	for _, pod := range pods {
		if pod.Status.NodeUuid == node.Metadata.Uuid {
			// Deduct the resources used by the running pods
			resources.Cpu -= pod.Spec.ResourceRequirement.Cpu
			resources.Ram -= pod.Spec.ResourceRequirement.Ram
			resources.Storage -= pod.Spec.ResourceRequirement.Storage
		}
	}
	return resources
}

func (ns *Scheduler) Run() {
	log.Printf("Starting scheduler")
	go ns.start()
}
