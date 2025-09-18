package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	clusterName string
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
	rootCmd.PersistentFlags().StringVar(&clusterName, "name", "kemu", "name of the KEMU cluster")
}
