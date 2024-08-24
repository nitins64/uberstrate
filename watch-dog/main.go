package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "gitbub.com/uberstrate/idl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func UpdateNodeTaint(nodeName string, tainted bool) {
	c, ctx := getClientWithContext()
	_, err := c.UpdateNodeTaint(ctx, &pb.UpdateNodeTaintRequest{Name: nodeName, Tainted: tainted})
	if err != nil {
		log.Fatalf("error calling function UpdateNodeTaint: %v", err)
	}
}

func getClientWithContext() (pb.StateStoreServiceClient, context.Context) {
	conn, err := getConnection()
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}

	c := pb.NewStateStoreServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return c, ctx
}

func main() {

	if len(os.Args) < 2 {
		log.Fatalf("expected 'nodeStore' subcommand")
	}

	switch os.Args[1] {
	case "node":
		switch os.Args[2] {
		case "tainted":
			nodeName := os.Args[3]
			UpdateNodeTaint(nodeName, true)
		case "untainted":
			nodeName := os.Args[3]
			UpdateNodeTaint(nodeName, false)
		default:
			log.Fatalf("unknown command: %s", os.Args[2])
		}
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}
