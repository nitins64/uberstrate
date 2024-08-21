package main

import (
	"fmt"
	pb "gitbub.com/uberstrate/idl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	port int32 = 50051
)

func getConnection() (*grpc.ClientConn, error) {
	target := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %v", target, err)
	}
	return conn, nil
}

// Define a custom error type
func getClientWithContext() pb.StateStoreServiceClient {
	conn, err := getConnection()
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}

	c := pb.NewStateStoreServiceClient(conn)
	return c
}

func main() {
	c := getClientWithContext()
	s := NewRepairEngine(c)
	s.Run()
	select {}
}
