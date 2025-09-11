package cmd

import (
	"os"

	"datastrophic.io/kemu/pkg/cluster"
	"github.com/spf13/cobra"
)

var (
	clusterConfigPath string
	kubeconfig        string
)

var clusterCmd = &cobra.Command{
	Use:   "create-cluster",
	Short: "Create KWOK cluster using provided configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.Flags().Parse(os.Args); err != nil {
			return err
		}

		err := cluster.CreateClusterWithConfig("kwok", clusterConfigPath, kubeconfig)
		if err != nil {
			return err
		}

		return cluster.CreateClusterNodes(clusterConfigPath, kubeconfig)
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.Flags().StringVar(&clusterConfigPath, "cluster-config", "", "KWOK cluster configuration file")
	clusterCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "KUBECONFIG file path for target KWOK cluster")
}
