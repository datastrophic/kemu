package cluster

import (
	"os"

	"datastrophic.io/kemu/pkg/api"
	"gopkg.in/yaml.v3"
)

func parseConfig(path string) (api.ClusterConfig, error) {
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
