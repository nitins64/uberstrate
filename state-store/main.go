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

func (s *server) SayHello(ctx context.Context, in *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	return &pb.HelloWorldResponse{Message: fmt.Sprintf("Hello, World to %s! ", in.Name),
		Model: &pb.Model{Name: "model"}}, nil
}

func (s *server) ReLoadNodes(ctx context.Context, in *pb.ReLoadNodeRequest) (*pb.ReLoadNodeResponse, error) {
	nodeStore := GetNodeStore()
	err := nodeStore.ReLoad(in.Path)
	if err != nil {
		return nil, err
	}
	return &pb.ReLoadNodeResponse{}, nil
}

func (s *server) PrintNodes(ctx context.Context, in *pb.PrintNodeRequest) (*pb.PrintNodeResponse, error) {
	nodeStore := GetNodeStore()
	err := nodeStore.PrintNodes()
	if err != nil {
		return nil, err
	}
	return &pb.PrintNodeResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStateStoreServiceServer(s, &server{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
