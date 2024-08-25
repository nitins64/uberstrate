package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "gitbub.com/uberstrate/idl"
	"google.golang.org/grpc"
)

type server struct {
	pb.StateStoreServiceServer
}

// Define a custom error type
type OperationNotAllowedError struct {
	Operation string
	Message   string
}

func (e *OperationNotAllowedError) Error() string {
	return fmt.Sprintf("operation not allowed: %s. %s", e.Operation, e.Message)
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	return &pb.HelloWorldResponse{Message: fmt.Sprintf("Hello, World to %s! ", in.Name),
		Model: &pb.Model{Name: "model"}}, nil
}

func (s *server) LoadNodes(ctx context.Context, in *pb.LoadNodeRequest) (*pb.LoadNodeResponse, error) {
	nodeStore := GetNodeStore()
	err := nodeStore.Load(in.Path)
	if err != nil {
		return nil, err
	}
	return &pb.LoadNodeResponse{}, nil
}

func (s *server) LoadPods(ctx context.Context, in *pb.LoadPodRequest) (*pb.LoadPodResponse, error) {
	ps := GetPodStore()
	err := ps.Load(in.Path)
	if err != nil {
		return nil, err
	}
	return &pb.LoadPodResponse{}, nil
}

func (s *server) PrintNodes(ctx context.Context, in *pb.PrintNodeRequest) (*pb.PrintNodeResponse, error) {
	nodeStore := GetNodeStore()
	err := nodeStore.PrintNodes()
	if err != nil {
		return nil, err
	}
	return &pb.PrintNodeResponse{}, nil
}

func (s *server) GetPods(ctx context.Context, in *pb.GetPodRequest) (*pb.GetPodResponse, error) {
	ps := GetPodStore()
	pods := ps.GetPods(in)
	return &pb.GetPodResponse{Pods: pods}, nil
}

func (s *server) UpdatePods(ctx context.Context, in *pb.UpdatePodRequest) (*pb.UpdatePodResponse, error) {
	ps := GetPodStore()
	err := ps.UpdatePods(in.Pods)
	if err != nil {
		return nil, err
	}
	return &pb.UpdatePodResponse{}, nil
}

func (s *server) UpdateNodeTainted(ctx context.Context, in *pb.UpdateNodeTaintRequest) (*pb.UpdateNodeTaintResponse, error) {
	nodeStore := GetNodeStore()
	err := nodeStore.UpdateNodeTaint(in)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateNodeTaintResponse{}, nil
}

func (s *server) GetNodes(ctx context.Context, in *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	nodeStore := GetNodeStore()
	nodes := nodeStore.GetNodes(in)
	return &pb.GetNodeResponse{Nodes: nodes}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	// Load the nodes and pods
	nodeStore := GetNodeStore()
	nodeStore.Load("../crane/nodes.yaml")
	ps := GetPodStore()
	ps.Load("../deployment/pods.yaml")

	s := grpc.NewServer()
	pb.RegisterStateStoreServiceServer(s, &server{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
