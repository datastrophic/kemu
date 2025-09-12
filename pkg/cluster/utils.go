package cluster

import (
	"log/slog"
	"os"

	"datastrophic.io/kemu/pkg/api"
	helmclient "github.com/mittwald/go-helm-client"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func kubeClientFromConfig(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func helmClientFromConfig(kubeconfig, namespace string) (helmclient.Client, error) {
	kubecfgBytes, err := os.ReadFile(kubeconfig)
	if err != nil {
		return nil, err
	}

	options := &helmclient.KubeConfClientOptions{
		Options: &helmclient.Options{
			Namespace: namespace,
			Debug:     true,
		},
		KubeConfig: kubecfgBytes,
	}

	slog.Info("initializing helm client", "kubeconfig", kubeconfig)
	return helmclient.NewClientFromKubeConf(options)
}

func parseKemuClusterConfig(path string) (api.ClusterConfig, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return api.ClusterConfig{}, err
	}

	var clusterConfig api.ClusterConfig

	if err := yaml.Unmarshal(f, &clusterConfig); err != nil {
		return api.ClusterConfig{}, err
	}

	return clusterConfig, nil
}

func writeTempFile(kindClusterConfig string) (*os.File, error) {
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
