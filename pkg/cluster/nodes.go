package cluster

import (
	"context"
	"fmt"
	"log/slog"

	"datastrophic.io/kemu/pkg/api"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
)

const (
	HostnameLabel      = "kubernetes.io/hostname"
	InstanceTypeLabel  = "node.kubernetes.io/instance-type"
	ZoneLabel          = "topology.kubernetes.io/zone"
	ManagedByKemuLabel = "kemu.datastrophic.io/managed"
)

func CreateClusterNodes(config api.ClusterConfig, kubeconfig string) error {
	slog.Info("creating KWOK cluster nodes")

	kubeClient, err := kubeClientFromConfig(kubeconfig)
	if err != nil {
		return err
	}

	var nodes []corev1.Node
	for _, nodeGroup := range config.Spec.NodeGroups {
		slog.Info("processing", "node group", nodeGroup.Name)
		for _, placement := range nodeGroup.Placement {
			nodes = append(nodes, createNodes(nodeGroup, placement)...)
		}
	}

	for _, node := range nodes {
		slog.Info("creating", "node", node.Name)
		_, err := kubeClient.CoreV1().Nodes().Create(context.TODO(), &node, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	slog.Info("nodes created", "count", len(nodes))
	return nil
}

func createNodes(nodeGroup api.NodeGroup, placement api.Placement) []corev1.Node {
	namePrefix := fmt.Sprintf("%s-%s", nodeGroup.Name, placement.AvailabilityZone)
	var nodes []corev1.Node

	for i := 0; i < placement.Replicas; i++ {
		hostname := fmt.Sprintf("%s-%d", namePrefix, i)

		annotations := map[string]string{
			"kwok.x-k8s.io/node": "fake",
		}

		labels := map[string]string{
			"kubernetes.io/arch": "arm64",
			"kubernetes.io/os":   "kemu",
			"kubernetes.io/role": "agent",
			"type":               "kwok",
			ManagedByKemuLabel:   "true",
			HostnameLabel:        hostname,
			InstanceTypeLabel:    nodeGroup.Name,
			ZoneLabel:            placement.AvailabilityZone,
		}

		for k, v := range nodeGroup.NodeTemplate.Labels {
			labels[k] = v
		}

		resources := make(map[corev1.ResourceName]resource.Quantity)
		for name, quantity := range nodeGroup.NodeTemplate.Capacity {
			resources[corev1.ResourceName(name)] = resource.MustParse(quantity)
		}

		node := corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:        hostname,
				Annotations: annotations,
				Labels:      labels,
			},
			Spec: corev1.NodeSpec{
				Taints: []corev1.Taint{
					{Key: "kwok.x-k8s.io/node", Effect: "NoSchedule", Value: "fake"},
				},
			},
			Status: corev1.NodeStatus{
				Capacity:    resources,
				Allocatable: resources,
				NodeInfo: corev1.NodeSystemInfo{
					Architecture:    "arm64",
					OperatingSystem: "kemu",
					KubeletVersion:  "fake",
					MachineID:       string(uuid.NewUUID()),
				},
				Phase: corev1.NodeRunning,
			},
		}

		nodes = append(nodes, node)
	}

	return nodes
}
