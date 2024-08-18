package main

import (
	"context"
	pb "gitbub.com/uberstrate/idl"
	"log"
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
	log.Printf("Scheduler loop")
	log.Printf("Getting nodes with min_generation: %d", 0)
	nodes, err := s.getNodes(0)
	if err != nil {
		log.Printf("error calling function GetNodes: %v", err)
	} else {
		log.Printf("Response from gRPC server's GetNodes total node: %d", len(nodes))
	}

	log.Printf("Getting pods with phase: %s", "RUNNING")
	pods, err := s.getPods(false, "PENDING_NODE_ASSIGNMENT")
	if err != nil {
		log.Fatalf("error calling function GetPods: %v", err)
	}
	log.Printf("Response from gRPC server's GetPods total pod: %d", len(pods))
}

func (ns *Scheduler) Run() {
	log.Printf("Starting scheduler")
	go ns.start()
}
