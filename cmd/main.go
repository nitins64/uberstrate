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

func HandleLoadNode(filePath string) {
	c, ctx := getClientWithContext()
	r, err := c.LoadNodes(ctx, &pb.LoadNodeRequest{Path: filePath})
	if err != nil {
		log.Fatalf("error calling function LoadNodes: %v", err)
	}

	log.Printf("Response from gRPC server's LoadNodes function: %+v", r)
}

func HandleLoadPods(filePath string) {
	c, ctx := getClientWithContext()

	r, err := c.LoadPods(ctx, &pb.LoadPodRequest{Path: filePath})
	if err != nil {
		log.Fatalf("error calling function LoadPods: %v", err)
	}
	log.Printf("Response from gRPC server's LoadPods function: %+v", r)
}

func PrintNodes() {
	c, ctx := getClientWithContext()
	c.PrintNodes(ctx, &pb.PrintNodeRequest{})
}

// Map to convert string to enum
var podPhaseMap = map[string]pb.PodPhase{
	"PENDING_NODE_ASSIGNMENT": pb.PodPhase_PENDING_NODE_ASSIGNMENT,
	"NODE_ASSIGNED":           pb.PodPhase_NODE_ASSIGNED,
	"RUNNING":                 pb.PodPhase_RUNNING,
	"FAILED":                  pb.PodPhase_FAILED,
}

func GetPods(all bool, phase string) {
	c, ctx := getClientWithContext()
	r, err := c.GetPods(ctx, &pb.GetPodRequest{All: all, Phase: podPhaseMap[phase]})
	if err != nil {
		log.Fatalf("error calling function GetPods: %v", err)
	}
	log.Printf("Response from gRPC server's GetPods pod: %v", r.Pods)
}

func GetNodes(minGeneration int64) {
	log.Printf("Getting nodes with min_generation: %d", minGeneration)
	c, ctx := getClientWithContext()
	r, err := c.GetNodes(ctx, &pb.GetNodeRequest{AboveGenerationNumber: minGeneration})
	if err != nil {
		log.Fatalf("error calling function GetNodes: %v", err)
	}
	log.Printf("Response from gRPC server's GetNodes %v", r.Nodes)
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
	case "nodeStore":
		switch os.Args[2] {
		case "load":
			load := flag.NewFlagSet("load", flag.ExitOnError)
			path := load.String("path", "", "path to the node store file")
			load.Parse(os.Args[3:]) // Adjusted to parse the correct arguments
			log.Printf("Loading nodes from %s", *path)
			HandleLoadNode(*path)
		case "print":
			// Handle print command
			PrintNodes()
		case "get":
			get := flag.NewFlagSet("get", flag.ExitOnError)
			min_generation := get.Int64("min_generation", 0, "")
			get.Parse(os.Args[3:])
			// Handle print command
			GetNodes(*min_generation)
		default:
			log.Fatalf("unknown command: %s", os.Args[2])
		}

	case "podStore":
		switch os.Args[2] {
		case "load":
			load := flag.NewFlagSet("load", flag.ExitOnError)
			path := load.String("path", "", "path to the node store file")
			load.Parse(os.Args[3:]) // Adjusted to parse the correct arguments
			log.Printf("Loading pods from %s", *path)
			HandleLoadPods(*path)
		case "get":
			get := flag.NewFlagSet("get", flag.ExitOnError)
			all := get.Bool("all", false, "")
			phase := get.String("phase", "", "")
			get.Parse(os.Args[3:])
			// Handle print command
			GetPods(*all, *phase)
		default:
			log.Fatalf("unknown command: %s", os.Args[2])
		}
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}
