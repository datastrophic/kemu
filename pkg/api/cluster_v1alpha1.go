package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterConfig struct {
	metav1.TypeMeta `yaml:",inline"`
	Spec            ClusterSpec `yaml:"spec"`
}

type ClusterSpec struct {
	NodeGroups    []NodeGroup    `yaml:"nodeGroups"`
	KindConfig    string         `yaml:"kindConfig"`
	ClusterAddons []ClusterAddon `yaml:"clusterAddons"`
}

type NodeGroup struct {
	Name         string       `yaml:"name"`
	Placement    []Placement  `yaml:"placement"`
	NodeTemplate NodeTemplate `yaml:"nodeTemplate"`
}

type Placement struct {
	AvailabilityZone string `yaml:"availabilityZone"`
	Replicas         int    `yaml:"replicas"`
}

type Resources map[string]string

type NodeTemplate struct {
	metav1.ObjectMeta `yaml:"metadata,omitempty"`
	Capacity          Resources `yaml:"capacity"`
}

type ClusterAddon struct {
	Name         string `yaml:"name"`
	Namespace    string `yaml:"namespace"`
	Chart        string `yaml:"chart"`
	RepoName     string `yaml:"repoName"`
	RepoURL      string `yaml:"repoURL"`
	Version      string `yaml:"version"`
	ValuesObject string `yaml:"valuesObject"`
}
