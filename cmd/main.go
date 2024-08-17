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

func HandleReloadNode(filePath string) {
	c, ctx := getClientWithContext()

	r, err := c.ReLoadNodes(ctx, &pb.ReLoadNodeRequest{Path: filePath})
	if err != nil {
		log.Fatalf("error calling function ReLoadNodes: %v", err)
	}

	log.Printf("Response from gRPC server's ReLoadNodes function: %v", r)
}

func PrintNodes() {
	c, ctx := getClientWithContext()
	c.PrintNodes(ctx, &pb.PrintNodeRequest{})
}

func getClientWithContext() (pb.StateStoreServiceClient, context.Context) {
	conn, err := getConnection()
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}

	c := pb.NewStateStoreServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	return c, ctx
}

func main() {

	if len(os.Args) < 2 {
		log.Fatalf("expected 'nodeStore' subcommand")
	}

	switch os.Args[1] {
	case "nodeStore":
		switch os.Args[2] {
		case "load":
			load := flag.NewFlagSet("load", flag.ExitOnError)
			path := load.String("path", "", "path to the node store file")
			load.Parse(os.Args[3:]) // Adjusted to parse the correct arguments
			log.Printf("Loading nodes from %s", *path)
			HandleReloadNode(*path)
		case "print":
			// Handle print command
			PrintNodes()
		default:
			log.Fatalf("unknown command: %s", os.Args[2])
		}

	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}
