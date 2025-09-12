package cluster

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"datastrophic.io/kemu/pkg/api"
	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/third_party/kind"
)

func createKindClusterWithConfig(config api.ClusterConfig, name, kubeconfig string) error {
	slog.Info("creating kind cluster", "name", name, "kubeconfig", kubeconfig)

	var configFile *os.File
	var err error
	if len(config.Spec.KindConfig) > 0 {
		configFile, err = writeTempFile(config.Spec.KindConfig)
		if err != nil {
			return err
		}
		defer os.Remove(configFile.Name())
	}

	// NOTE: there doesn't seem to be a way of passing the KUBECONFIG location
	// to the Kind cluster provider, but it respects the env var setting.
	err = os.Setenv("KUBECONFIG", kubeconfig)

	kindClusterProvider := kind.NewProvider().SetDefaults().WithName(name)
	if configFile != nil {
		slog.Info(fmt.Sprintf("using provided kind config:\n%s", config.Spec.KindConfig))
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
	err = kindClusterProvider.WaitForControlPlane(context.TODO(), client)
	if err != nil {
		return err
	}
	slog.Info("kind control plane is ready")
	return nil
}

func deleteKindCluster(name string) error {
	slog.Info("deleting kind cluster", "name", name)
	return kind.NewProvider().SetDefaults().WithName(name).Destroy(context.Background())
}
