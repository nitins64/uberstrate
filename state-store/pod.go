package main

import (
	"log"
	"os"
	"reflect"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gitbub.com/uberstrate/idl"
	"gopkg.in/yaml.v2"
)

// NodeSelection Define the NodeSelection struct
type NodeSelection struct {
	Region   string `yaml:"region"`
	Zone     string `yaml:"zone"`
	Cluster  string `yaml:"cluster"`
	OS       string `yaml:"os"`
	DiskType string `yaml:"disk_type"`
}

// Resource Define the Resource struct
type Resource struct {
	CPU     int `yaml:"cpu"`
	RAM     int `yaml:"ram"`
	Storage int `yaml:"storage"`
}

// Pod Define the Pod struct
type Pod struct {
	Name          string            `yaml:"name"`
	NodeSelection NodeSelection     `yaml:"nodeSelection"`
	Resource      Resource          `yaml:"resource"`
	Priority      int               `yaml:"priority"`
	Image         string            `yaml:"image"`
	Labels        map[string]string `yaml:"labels"`
}

type PodStore struct {
	NameToPod         map[string]Pod
	NameToPodProto    map[string]*pb.Pod
	DeploymentPodPath string
	Generation        int
	loaded            bool
	mutex             sync.Mutex
}

var (
	instancePodStore *PodStore
	onceInitPodStore sync.Once
)

func (ps *PodStore) loadInternal() error {
	pods, err := ps.loadPods()
	if err != nil {
		return err
	}
	//log.Printf("Reloaded pods from %s", ps.DeploymentPodPath)
	for _, pod := range pods {
		if _, exists := ps.NameToPod[pod.Name]; exists {
			if reflect.DeepEqual(ps.NameToPod[pod.Name], pod) {
				continue
			}
			log.Printf("Pod %s definition has changed. "+
				"Pods is immutable can't be changed. Ignoring the changes", pod.Name)
			continue
		} else {
			log.Printf("Adding new pod %s.", pod.Name)
		}
		ps.NameToPod[pod.Name] = pod
		ps.NameToPodProto[pod.Name] = CreateProtoForPod(ps, pod)
	}
	ps.loaded = true
	return nil
}

func CreateProtoForPod(ps *PodStore, pod Pod) *pb.Pod {
	// Implementation of CreateProtoForPod
	return &pb.Pod{
		Ot: &pb.ObjectType{
			Version: "v1",
			Kind:    "Pod",
		},
		Metadata: &pb.Metadata{
			Name:       pod.Name,
			Uuid:       pod.Name,
			CreateTime: timestamppb.Now(),
			// TODO: Implement this function
			// Not used since POD spec is immutable.
			GenerationNumber: 0,
			Labels:           pod.Labels,
		},
		Spec: &pb.PodSpec{
			Image: pod.Image,
			ResourceRequirement: &pb.Resource{
				Cpu:     int64(pod.Resource.CPU),
				Ram:     int64(pod.Resource.RAM),
				Storage: int64(pod.Resource.Storage),
			},
			Priority: int64(pod.Priority),
			NodeSelectorLabels: map[string]string{
				"region":    pod.NodeSelection.Region,
				"zone":      pod.NodeSelection.Zone,
				"cluster":   pod.NodeSelection.Cluster,
				"os":        pod.NodeSelection.OS,
				"disk_type": pod.NodeSelection.DiskType,
			},
		},
		Status: &pb.PodStatus{
			Phase: pb.PodPhase_PENDING_NODE_ASSIGNMENT,
		},
	}
}

func (ps *PodStore) scheduleLoad() {
	uptimeTicker := time.NewTicker(5 * time.Second)
	for range uptimeTicker.C {
		ps.mutex.Lock()
		ps.loadInternal()
		ps.mutex.Unlock()
	}
}

func (ps *PodStore) Load(filePath string) error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	if ps.loaded {
		return &OperationNotAllowedError{
			Operation: "Load",
			Message:   "Load already done and running in background"}
	}
	ps.DeploymentPodPath = filePath
	err := ps.loadInternal()
	go ps.scheduleLoad()
	return err
}

func GetPodStore() *PodStore {
	onceInitPodStore.Do(func() {
		instancePodStore = &PodStore{
			NameToPod:      make(map[string]Pod),
			NameToPodProto: make(map[string]*pb.Pod),
		}
	})
	return instancePodStore
}

func (ps *PodStore) loadPods() ([]Pod, error) {
	// Read the YAML file
	data, err := os.ReadFile(ps.DeploymentPodPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Unmarshal the YAML data into the Pod struct
	var pods []Pod
	err = yaml.Unmarshal(data, &pods)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return pods, nil
}

func (ps *PodStore) GetPods(in *pb.GetPodRequest) (pods []*pb.Pod) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	//log.Printf("Number of pods: %d", len(ps.NameToPodProto))

	pods = make([]*pb.Pod, 0, len(ps.NameToPodProto))
	for _, pod := range ps.NameToPodProto {
		if in.All || pod.Status.Phase == in.Phase {
			pods = append(pods, pod)
		}
	}
	//log.Printf("Total pods: %d", len(pods))
	return pods
}

func (ps *PodStore) UpdatePods(pod []*pb.Pod) error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	// TODO: Need concurrent control here or any place where we are updating state store.
	for _, pod := range pod {
		if _, exists := ps.NameToPod[pod.Metadata.Name]; !exists {
			return &OperationNotAllowedError{
				Operation: "UpdatePod",
				Message:   "Pod does not exist"}
		}
		ps.NameToPodProto[pod.Metadata.Name] = pod
		log.Printf("Updated pod name:%s status:%s condition:%s assinged-node:%s", pod.Metadata.Name,
			pod.Status.Phase, pod.Status.Condition, pod.Status.NodeUuid)
	}
	return nil
}
