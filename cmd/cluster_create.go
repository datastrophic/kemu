package cmd

import (
	"fmt"
	"os"

	"datastrophic.io/kemu/pkg/cluster"
	"github.com/spf13/cobra"
)

var (
	clusterConfig string
	kubeconfig        string
)

var createClusterCmd = &cobra.Command{
	Use:   "create-cluster",
	Short: "Create KEMU cluster using provided configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.Flags().Parse(os.Args); err != nil {
			return err
		}
		if len(clusterConfig) == 0 {
			return fmt.Errorf("--cluster-config is required")
		}

		return cluster.CreateKemuCluster(clusterConfig, clusterName, kubeconfig)
	},
}

func init() {
	rootCmd.AddCommand(createClusterCmd)
	createClusterCmd.Flags().StringVar(&clusterConfig, "cluster-config", "", "KEMU cluster configuration file or URL")
	createClusterCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "kemu.config", "KUBECONFIG file for accessing created KEMU cluster")
}
