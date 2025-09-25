package test

import (
	"context"
	"fmt"
	"os/exec"

	"datastrophic.io/kemu/pkg/cluster"
	"datastrophic.io/kemu/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Create client from kubeconfig with assertions.
func getClient(kubeconfig string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	Expect(err).NotTo(HaveOccurred(), "failed to build config from kubeconfig")
	c, err := kubernetes.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred(), "failed to build client from kubeconfig")
	return c
}

// Delete cluster with assertions.
func deleteCluster(clusterName string) {
	err := cluster.DeleteKemuCluster(clusterName)
	Expect(err).NotTo(HaveOccurred(), "failed to delete cluster")

	cmd := exec.Command("kind", "get", "clusters")
	output, err := utils.Run(cmd)
	Expect(err).NotTo(HaveOccurred(), "failed to list kind clusters")
	Expect(output).NotTo(ContainSubstring(clusterName), "kind cluster e2e-simple still exist")
}

var _ = Describe("kemu API", Ordered, func() {
	clusterName := "it-simple"
	kubeconfig := ".run/it-simple.config"

	Context("CreateKemuCluster", func() {
		It("should create a simple Kind cluster based on cluster spec", func() {
			err := cluster.CreateKemuCluster("test/testdata/simple.yaml", clusterName, kubeconfig)
			Expect(err).NotTo(HaveOccurred(), "failed to create cluster")

			client := getClient(kubeconfig)
			nodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred(), "failed to list nodes")
			Expect(len(nodes.Items)).To(Equal(1), fmt.Sprintf("expected 1 node, got %d", len(nodes.Items)))
		})
	})
	Context("DeleteKemuCluster", func() {
		It("should delete created Kind cluster", func() {
			deleteCluster(clusterName)
		})
	})
})

var _ = Describe("kemu API", Ordered, func() {
	clusterName := "it-with-kind-config"
	kubeconfig := ".run/it-with-kind.config"

	Context("CreateKemuCluster", func() {
		It("should create a Kind cluster with custom Kind config", func() {
			err := cluster.CreateKemuCluster("test/testdata/with-kind-config.yaml", clusterName, kubeconfig)
			Expect(err).NotTo(HaveOccurred(), "failed to create cluster")

			client := getClient(kubeconfig)
			nodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
				LabelSelector: "kemu.datastrophic.io/e2e=true",
			})
			Expect(err).NotTo(HaveOccurred(), "failed to list nodes")
			Expect(len(nodes.Items)).To(Equal(3), fmt.Sprintf("expected 3 nodes, got %d", len(nodes.Items)))
		})
	})
	Context("DeleteKemuCluster", func() {
		It("should delete created Kind cluster with custom Kind config", func() {
			deleteCluster(clusterName)
		})
	})
})

var _ = Describe("kemu API", Ordered, func() {
	clusterName := "it-with-addons"
	kubeconfig := ".run/it-with-addons.config"

	Context("CreateKemuCluster", func() {
		It("should create a Kind cluster with addons", func() {
			err := cluster.CreateKemuCluster("test/testdata/with-addons.yaml", clusterName, kubeconfig)
			Expect(err).NotTo(HaveOccurred(), "failed to create cluster")

			client := getClient(kubeconfig)
			d, err := client.AppsV1().Deployments("kube-system").Get(context.Background(), "kwok-controller", metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred(), "failed to get addon deployment")
			Expect(d.Status.ReadyReplicas).To(Equal(int32(1)), fmt.Sprintf("expected kwok addon deployment to have 1 ready replica but got %d", d.Status.ReadyReplicas))
		})
	})
	Context("DeleteKemuCluster", func() {
		It("should delete created Kind cluster with addons", func() {
			deleteCluster(clusterName)
		})
	})
})

var _ = Describe("kemu API", Ordered, func() {
	clusterName := "it-with-kwok-nodes"
	kubeconfig := ".run/it-with-kwok.config"

	Context("CreateKemuCluster", func() {
		It("should create a Kind cluster with kwok nodes", func() {
			err := cluster.CreateKemuCluster("test/testdata/with-kwok-nodes.yaml", clusterName, kubeconfig)
			Expect(err).NotTo(HaveOccurred(), "failed to create cluster")

			client := getClient(kubeconfig)
			data := []struct {
				instanceType string
				zone         string
				expected     int
			}{
				{
					"a2-ultragpu-8g",
					"use1",
					5,
				},
				{
					"a2-ultragpu-8g",
					"use2",
					5,
				},
				{
					"a2-ultragpu-8g",
					"use3",
					5,
				},
				{
					"a3-highgpu-8g",
					"use1",
					10,
				},
				{
					"a3-highgpu-8g",
					"use2",
					10,
				},
			}
			for _, tc := range data {
				nodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
					LabelSelector: fmt.Sprintf("node.kubernetes.io/instance-type=%s,topology.kubernetes.io/zone=%s", tc.instanceType, tc.zone),
				})
				Expect(err).NotTo(HaveOccurred(), "failed to list nodes")
				Expect(len(nodes.Items)).To(Equal(tc.expected), fmt.Sprintf("expected %d nodes, got %d", tc.expected, len(nodes.Items)))
			}
		})
	})
	Context("DeleteKemuCluster", func() {
		It("should delete created Kind cluster with kwok nodes", func() {
			deleteCluster(clusterName)
		})
	})
})

var _ = Describe("kemu API", Ordered, func() {
	clusterName := "it-with-full-config"
	kubeconfig := ".run/it-with-full-config.config"

	Context("CreateKemuCluster", func() {
		It("should create a Kind cluster with full config", func() {
			err := cluster.CreateKemuCluster("test/testdata/with-full-config.yaml", clusterName, kubeconfig)
			Expect(err).NotTo(HaveOccurred(), "failed to create cluster")

			client := getClient(kubeconfig)
			d, err := client.AppsV1().Deployments("kube-system").Get(context.Background(), "kwok-controller", metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred(), "failed to get addon deployment")
			Expect(d.Status.ReadyReplicas).To(Equal(int32(1)), fmt.Sprintf("expected kwok addon deployment to have 1 ready replica but got %d", d.Status.ReadyReplicas))

			d, err = client.AppsV1().Deployments("monitoring").Get(context.Background(), "prometheus-grafana", metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred(), "failed to get addon deployment")
			Expect(d.Status.ReadyReplicas).To(Equal(int32(1)), fmt.Sprintf("expected prometheus addon deployment to have 1 ready replica but got %d", d.Status.ReadyReplicas))

			data := []struct {
				instanceType string
				zone         string
				expected     int
			}{
				{
					"a2-ultragpu-8g",
					"use1",
					5,
				},
				{
					"a2-ultragpu-8g",
					"use2",
					5,
				},
				{
					"a2-ultragpu-8g",
					"use3",
					5,
				},
				{
					"a3-highgpu-8g",
					"use1",
					5,
				},
				{
					"a3-highgpu-8g",
					"use2",
					5,
				},
				{
					"a3-ultragpu-8g",
					"use1",
					5,
				},
			}
			for _, tc := range data {
				nodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
					LabelSelector: fmt.Sprintf("node.kubernetes.io/instance-type=%s,topology.kubernetes.io/zone=%s", tc.instanceType, tc.zone),
				})
				Expect(err).NotTo(HaveOccurred(), "failed to list nodes")
				Expect(len(nodes.Items)).To(Equal(tc.expected), fmt.Sprintf("expected %d nodes, got %d", tc.expected, len(nodes.Items)))
			}
		})
	})
	Context("DeleteKemuCluster", func() {
		It("should delete created Kind cluster with full config", func() {
			deleteCluster(clusterName)
		})
	})
})
