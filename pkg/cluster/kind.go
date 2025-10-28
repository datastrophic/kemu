package cluster

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/pkg/utils"
	"sigs.k8s.io/e2e-framework/third_party/kind"
)

func kindClusterExists(name string) bool {
	clusters := utils.FetchCommandOutput("kind get clusters")
	for _, c := range strings.Split(clusters, "\n") {
		if c == name {
			return true
		}
	}
	return false
}

func createKindClusterWithConfig(kindConfig string, name, kubeconfig string) error {
	slog.Info("creating kind cluster", "name", name, "kubeconfig", kubeconfig)

	var configFile *os.File
	var err error
	if len(kindConfig) > 0 {
		configFile, err = writeTempFile(kindConfig)
		if err != nil {
			return err
		}
		defer os.Remove(configFile.Name())
	}

	kindClusterProvider := kind.NewProvider().SetDefaults().WithName(name)
	args := []string{"--kubeconfig", kubeconfig}
	if configFile != nil {
		slog.Info(fmt.Sprintf("using provided kind config:\n%s", kindConfig))
		args = append(args, "--config", configFile.Name())
	}

	_, err = kindClusterProvider.Create(context.Background(), args...)
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
