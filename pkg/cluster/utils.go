package cluster

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
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

func parseKemuClusterConfig(loc string) (api.ClusterConfig, error) {
	var body []byte
	var err error

	u, err := url.Parse(loc)
	if err == nil && u.Scheme != "" && u.Host != "" {
		body, err = configFromURL(loc)
	} else {
		body, err = os.ReadFile(loc)
	}

	if err != nil {
		return api.ClusterConfig{}, err
	}

	var clusterConfig api.ClusterConfig

	if err := yaml.Unmarshal(body, &clusterConfig); err != nil {
		return api.ClusterConfig{}, err
	}

	return clusterConfig, nil
}

func configFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
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
