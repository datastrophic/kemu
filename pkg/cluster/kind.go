package cluster

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/third_party/kind"
)

func CreateClusterWithConfig(clusterName, configPath, kubeconfig string) error {
	slog.Info("creating kind cluster", "name", clusterName, "kubeconfig", kubeconfig)
	clusterConfig, err := parseConfig(configPath)
	if err != nil {
		return err
	}

	var configFile *os.File
	if len(clusterConfig.KindConfig) > 0 {
		configFile, err = writeConfig(clusterConfig.KindConfig)
		if err != nil {
			return err
		}
		defer os.Remove(configFile.Name())
	}

	// NOTE: there doesn't seem to be a way of passing the KUBECONFIG location
	// to the Kind cluster provider, but it respects the env var setting.
	err = os.Setenv("KUBECONFIG", kubeconfig)

	kindClusterProvider := kind.NewProvider().SetDefaults().WithName(clusterName)
	if configFile != nil {
		slog.Info(fmt.Sprintf("using provided kind config:\n%s", clusterConfig.KindConfig))
		_, err = kindClusterProvider.CreateWithConfig(context.TODO(), configFile.Name())
	} else {
		_, err = kindClusterProvider.Create(context.TODO())
	}
	if err != nil {
		return err
	}

	client, err := klient.NewWithKubeConfigFile(kubeconfig)
	if err != nil {
		return err
	}

	slog.Info("waiting for control plane to become ready")
	return kindClusterProvider.WaitForControlPlane(context.TODO(), client)
}

func writeConfig(kindClusterConfig string) (*os.File, error) {
	configFile, err := os.CreateTemp("", "kind.yaml")
	if err != nil {
		return nil, err
	}

	_, err = configFile.WriteString(kindClusterConfig)
	if err != nil {
		return nil, err
	}
	return configFile, nil
}
