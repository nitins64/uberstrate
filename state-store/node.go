package main

import (
	"log"
	"os"
	"reflect"
	"sync"

	pb "gitbub.com/uberstrate/idl"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/yaml.v2"
)

type Node struct {
	Name     string `yaml:"name"`
	Region   string `yaml:"region"`
	Zone     string `yaml:"zone"`
	Cluster  string `yaml:"cluster"`
	OS       string `yaml:"os"`
	DiskType string `yaml:"disk_type"`
	CPU      int    `yaml:"cpu"`
	RAM      int    `yaml:"ram"`
	Storage  int    `yaml:"storage"`
	Offline  bool   `yaml:"offline"`
}

type NodeStore struct {
	NameToNode      map[string]Node
	NameToNodeProto map[string]*pb.Node

	Generation int
	mutex      sync.Mutex
}

var (
	instance *NodeStore
	once     sync.Once
)

func GetNodeStore() *NodeStore {
	once.Do(func() {
		instance = &NodeStore{
			NameToNode:      make(map[string]Node),
			NameToNodeProto: make(map[string]*pb.Node),
		}
	})
	return instance
}

func (ns *NodeStore) PrintNodes() error {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	for _, node := range ns.NameToNode {
		log.Printf("Node: %+v", node)
	}
	log.Printf("Generation: %d", ns.Generation)
	return nil
}

func (ns *NodeStore) GetNodes(in *pb.GetNodeRequest) (nodes []*pb.Node) {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()
	log.Printf("Get nodes with AboveGenerationNumber: %d", in.AboveGenerationNumber)
	log.Printf("Number of nodes: %d", len(ns.NameToNodeProto))

	nodes = make([]*pb.Node, 0, len(ns.NameToNodeProto))
	for _, node := range ns.NameToNodeProto {
		if node.Metadata.GenerationNumber > in.AboveGenerationNumber {
			nodes = append(nodes, node)
		}
	}
	log.Printf("Total nodes: %d", len(nodes))
	return nodes
}

func CreateProtoForNode(ns *NodeStore, node Node) *pb.Node {
	tainted := ""
	if node.Offline {
		tainted = "NoSchedule"
	}

	return &pb.Node{
		Ot: &pb.ObjectType{
			Version: "v1",
			Kind:    "Node",
		},
		Metadata: &pb.Metadata{
			Name:             node.Name,
			Uuid:             node.Name,
			CreateTime:       timestamppb.Now(),
			GenerationNumber: int64(ns.Generation),
			Labels: map[string]string{
				"region":    node.Region,
				"zone":      node.Zone,
				"cluster":   node.Cluster,
				"os":        node.OS,
				"disk_type": node.DiskType,
			},
		},
		Spec: &pb.NodeSpec{
			Taint: tainted,
		},
		Status: &pb.NodeStatus{
			Capacity: &pb.Resource{
				Cpu:     int64(node.CPU),
				Ram:     int64(node.RAM),
				Storage: int64(node.Storage),
			},
			Phase: "Ready",
		},
	}
}

func (ns *NodeStore) ReLoad(filePath string) error {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	nodes, err := loadNodes(filePath)
	if err != nil {
		return err
	}
	ns.Generation++
	for _, node := range nodes {
		if _, exists := ns.NameToNode[node.Name]; exists {
			log.Printf("Node %s already exists.", node.Name)
			if reflect.DeepEqual(ns.NameToNode[node.Name], node) {
				continue
			}
			log.Printf("Node %s has changed, updating.", node.Name)
		} else {
			log.Printf("Adding new node %s.", node.Name)
		}
		ns.NameToNode[node.Name] = node
		ns.NameToNodeProto[node.Name] = CreateProtoForNode(ns, node)
	}
	return nil
}

func loadNodes(filePath string) (node []Node, err error) {
	// Read the YAML file
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	// Unmarshal the YAML data into a slice of Node structs
	var nodes []Node
	err = yaml.Unmarshal(yamlFile, &nodes)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML data: %v", err)
		return nil, err
	}
	return nodes, nil
}
