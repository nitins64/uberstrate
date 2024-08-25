package main

import (
	"context"
	"flag"
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

func UpdateNodeTaint(nodeName string, tainted bool, reason string) {
	c, ctx := getClientWithContext()
	nodes := []*pb.NodeTainted{&pb.NodeTainted{Name: nodeName, Tainted: tainted, Reason: reason}}
	_, err := c.UpdateNodeTainted(ctx, &pb.UpdateNodeTaintRequest{Nodes: nodes})
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
			tainted := flag.NewFlagSet("tainted", flag.ExitOnError)
			name := tainted.String("name", "", "Name of Node")
			reason := tainted.String("reason", "", "Reason")
			tainted.Parse(os.Args[3:]) // Adjusted to parse the correct arguments
			if *name == "" {
				log.Fatalf("name of node is required")
			}

			UpdateNodeTaint(*name, true, *reason)

		case "untainted":
			untainted := flag.NewFlagSet("untainted", flag.ExitOnError)
			name := untainted.String("name", "", "Name of Node")
			untainted.Parse(os.Args[3:]) // Adjusted to parse the correct arguments
			if name == nil || *name == "" {
				log.Fatalf("name Of node is required")
			}
			UpdateNodeTaint(*name, false, "")
		default:
			log.Fatalf("unknown command: %s", os.Args[2])
		}
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}
