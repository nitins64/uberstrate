package main

import (
	"context"
	pb "gitbub.com/uberstrate/idl"
	"log"
	"time"
)

type NodeAgent struct {
	client pb.StateStoreServiceClient
}

func NewNodeAgent(client pb.StateStoreServiceClient) *NodeAgent {
	return &NodeAgent{
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

func (n *NodeAgent) getPods(all bool, phase string) ([]*pb.Pod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := n.client.GetPods(ctx, &pb.GetPodRequest{All: all, Phase: podPhaseMap[phase]})
	if err != nil {
		return nil, err
	}
	return r.Pods, nil
}

func (n *NodeAgent) getNodes(minGeneration int64) ([]*pb.Node, error) {
	log.Printf("Getting nodes with min_generation: %d", minGeneration)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := n.client.GetNodes(ctx, &pb.GetNodeRequest{AboveGenerationNumber: minGeneration})
	if err != nil {
		return nil, err
	}
	return r.Nodes, nil
}

func (n *NodeAgent) start() {
	uptimeTicker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-uptimeTicker.C:
			n.loop()
		}
	}
}

func findFirst[T any](slice []T, condition func(T) bool) *T {
	for _, item := range slice {
		if condition(item) {
			return &item
		}
	}
	return nil
}

// For now node-agent, monitors all the nodes.
// In prod, we will have one node-agent deployed per node. And will monitor only that node.
func (n *NodeAgent) loop() {
	nodes, err := n.getNodes(0)
	if err != nil {
		log.Printf("error calling function GetNodes: %v", err)
	} else {
		log.Printf("Response from gRPC server'n GetNodes total node: %d", len(nodes))
	}

	log.Printf("Getting all podAlls")
	podAlls, err := n.getPods(true /* all */, "" /* phase */)
	if err != nil {
		log.Printf("error calling function GetPods: %v", err)
		return
	}
	log.Printf("Response from gRPC server'n GetPods total pod: %d", len(podAlls))

	availableResourcesOnNode := make(map[string]*pb.Resource)
	for _, node := range nodes {
		availableResourcesOnNode[node.Metadata.Uuid] = availableResources(node, podAlls)
	}

	// Check all the pod that need to run on the node
	// Change their status from NODE_ASSIGNED to RUNNING
	// Before running check whether the node has enough resources to run the pod.
	// If not, change the status to FAILED and PodCondition to REALLOCATION_REQUIRED
	// If the node has enough resources, change the status to RUNNING
	for _, pod := range podAlls {
		if pod.Status.Phase == pb.PodPhase_NODE_ASSIGNED {
			condition := func(node *pb.Node) bool {
				return pod.Status.NodeUuid == node.Metadata.Uuid
			}
			// Call the function with a slice of pods and the condition
			firstNode := findFirst(nodes, condition)
			if firstNode == nil {
				log.Println("No node found that meets the condition. Something went wrong")
				pod.Status.Phase = pb.PodPhase_FAILED
				pod.Status.Condition = pb.PodCondition_REALLOCATION_REQUIRED
			} else {
				node := *firstNode
				availableResources := availableResourcesOnNode[node.Metadata.Uuid]
				if availableResources.Cpu < pod.Spec.ResourceRequirement.Cpu ||
					availableResources.Ram < pod.Spec.ResourceRequirement.Ram ||
					availableResources.Storage < pod.Spec.ResourceRequirement.Storage {
					log.Printf("Node: %s does not have enough resources to run pod: %s", node.Metadata.Name, pod.Metadata.Name)
					pod.Status.Phase = pb.PodPhase_FAILED
					pod.Status.Condition = pb.PodCondition_REALLOCATION_REQUIRED
				} else {
					// Subtract the resources used by the pod from the availableResourcesOnNode
					availableResourcesOnNode[node.Metadata.Uuid].Cpu -= pod.Spec.ResourceRequirement.Cpu
					availableResourcesOnNode[node.Metadata.Uuid].Ram -= pod.Spec.ResourceRequirement.Ram
					availableResourcesOnNode[node.Metadata.Uuid].Storage -= pod.Spec.ResourceRequirement.Storage
					log.Printf("Node: %s has enough resources to run pod: %s", node.Metadata.Name, pod.Metadata.Name)
					pod.Status.Phase = pb.PodPhase_RUNNING
					pod.Status.Condition = pb.PodCondition_READY
				}
			}
			log.Printf("Updating pod: %s", pod.Metadata.Name)
			_, err = n.client.UpdatePods(context.Background(), &pb.UpdatePodRequest{Pods: []*pb.Pod{pod}})
			if err != nil {
				log.Printf("error calling function UpdatePods: %v", err)
			}
		}
	}

	// Now reassign pod that are on Nodes that are tainted.
	// Change their condition to REALLOCATION_REQUIRED
	for _, pod := range podAlls {
		if pod.Status.Phase != pb.PodPhase_PENDING_NODE_ASSIGNMENT &&
			pod.Status.Condition != pb.PodCondition_REALLOCATION_REQUIRED {
			condition := func(node *pb.Node) bool {
				return pod.Status.NodeUuid == node.Metadata.Uuid && node.Spec.Taint != ""
			}
			firstNode := findFirst(nodes, condition)
			if firstNode != nil {
				node := *firstNode
				log.Printf("Node: %s is tainted: %s. Need to reallocate pod: %s",
					node.Metadata.Name, node.Spec.Taint, pod.Metadata.Name)
				pod.Status.Condition = pb.PodCondition_REALLOCATION_REQUIRED
				log.Printf("Pod: %s condition changed to REALLOCATION_REQUIRED since node was tainted", pod.Metadata.Name)
				_, err = n.client.UpdatePods(context.Background(), &pb.UpdatePodRequest{Pods: []*pb.Pod{pod}})
				if err != nil {
					log.Printf("error calling function UpdatePods: %v", err)
				}
			}
		}
	}

	// Now check if Node is deleted from the state store and change the status of the pod to FAILED
	//and condition to REALLOCATION_REQUIRED
	for _, pod := range podAlls {
		condition := func(node *pb.Node) bool {
			return pod.Status.NodeUuid == node.Metadata.Uuid
		}
		firstNode := findFirst(nodes, condition)
		if firstNode == nil {
			log.Printf("Node: %s is deleted. Need to reallocate pod: %s", pod.Status.NodeUuid, pod.Metadata.Name)
			pod.Status.Condition = pb.PodCondition_REALLOCATION_REQUIRED
			pod.Status.Phase = pb.PodPhase_FAILED
			_, err = n.client.UpdatePods(context.Background(), &pb.UpdatePodRequest{Pods: []*pb.Pod{pod}})
			if err != nil {
				log.Printf("error calling function UpdatePods: %v", err)
			}
		}
	}
}

func availableResources(node *pb.Node, pods []*pb.Pod) *pb.Resource {
	if node.Spec.Taint != "" {
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

func (n *NodeAgent) Run() {
	log.Printf("Starting node agent")
	go n.start()
}
