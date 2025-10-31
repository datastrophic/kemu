package test

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/datastrophic/kemu/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/e2e-framework/support/kind"
)

var _ = Describe("kemu CLI", Ordered, func() {
	clusterName := "e2e-cli-test-cluster"
	kubeconfig := ".run/e2e-simple.conf"

	BeforeAll(func() {
		cmd := exec.Command("go", "build", "-o", ".run/kemu", "main.go")
		_, err := utils.Run(cmd)
		ExpectWithOffset(1, err).NotTo(HaveOccurred(), "failed to build the kemu binary")
	})

	AfterAll(func() {
		err := kind.NewProvider().SetDefaults().WithName(clusterName).Destroy(context.Background())
		Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("failed to destroy cluster %s", clusterName))
	})

	Context("create-cluster command", func() {
		It("should create an empty Kind cluster based on cluster spec", func() {
			cmd := exec.Command(".run/kemu", "create-cluster", "--cluster-config", "test/testdata/simple.yaml", "--name", clusterName, "--kubeconfig", kubeconfig)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "failed to create cluster")
			Expect(kubeconfig).To(BeAnExistingFile(), "expected kubeconfig file doesn't exist")

			config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
			Expect(err).NotTo(HaveOccurred(), "failed to build config from kubeconfig")
			client, err := kubernetes.NewForConfig(config)
			Expect(err).NotTo(HaveOccurred(), "failed to build client from kubeconfig")

			nodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred(), "failed to list nodes")
			Expect(len(nodes.Items)).To(Equal(1), fmt.Sprintf("expected 1 node, got %d", len(nodes.Items)))
		})
	})
	Context("delete-cluster command", func() {
		It("should delete an empty Kind cluster", func() {
			cmd := exec.Command(".run/kemu", "delete-cluster", "--name", clusterName)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "failed to delete cluster")

			cmd = exec.Command("kind", "get", "clusters")
			output, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "failed to list kind clusters")
			Expect(output).NotTo(ContainSubstring(clusterName), "kind cluster e2e-simple still exist")
		})
	})
})
