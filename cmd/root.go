package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	clusterName string
	kubeconfig  string
)

var rootCmd = &cobra.Command{
	Use:   "kemu",
	Short: "A Kubernetes Cluster emulation tool based on Kind and KWOK",
	Long: `Easily create emulated Kubernetes clusters with 1000s of nodes
for workload scheduling experimentation and analysis.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	createClusterCmd.PersistentFlags().StringVar(&clusterName, "name", "kemu", "name of the KEMU cluster")
	createClusterCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "KUBECONFIG file path for target KEMU cluster")
}
