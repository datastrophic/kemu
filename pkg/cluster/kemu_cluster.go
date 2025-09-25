package cluster

import (
	"fmt"
	"log/slog"
)

func CreateKemuCluster(configPath, name, kubeconfig string) error {
	slog.Info("creating KEMU cluster", "name", name)
	clusterConfig, err := parseKemuClusterConfig(configPath)
	if err != nil {
		return err
	}

	if kindClusterExists(name) {
		return fmt.Errorf("underlying kind cluster %q already exists. it needs to be deleted first", name)
	}
	if err = createKindClusterWithConfig(clusterConfig, name, kubeconfig); err != nil {
		return err
	}

	err = InstallOrUpgradeAddons(clusterConfig, kubeconfig)
	if err != nil {
		return err
	}
	err = CreateClusterNodes(clusterConfig, kubeconfig)
	if err != nil {
		return err
	}
	slog.Info("KEMU cluster created", "name", name)
	return nil
}

func DeleteKemuCluster(name string) error {
	slog.Info("deleting KEMU cluster", "name", name)
	if err := deleteKindCluster(name); err != nil {
		return err
	}
	slog.Info("KEMU cluster deleted", "name", name)
	return nil
}
