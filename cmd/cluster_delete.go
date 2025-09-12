package cmd

import (
	"os"

	"datastrophic.io/kemu/pkg/cluster"
	"github.com/spf13/cobra"
)

var deleteClusterCmd = &cobra.Command{
	Use:   "delete-cluster",
	Short: "Delete KEMU cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.Flags().Parse(os.Args); err != nil {
			return err
		}

		return cluster.DeleteKemuCluster(clusterName)
	},
}

func init() {
	rootCmd.AddCommand(deleteClusterCmd)
}
