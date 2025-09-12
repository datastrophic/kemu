package cluster

import (
	"context"
	"log/slog"
	"time"

	"datastrophic.io/kemu/pkg/api"
	"helm.sh/helm/v3/pkg/repo"

	helmclient "github.com/mittwald/go-helm-client"
)

func InstallOrUpgradeAddons(config api.ClusterConfig, kubeconfig string) error {
	slog.Info("installing cluster addons")
	helmClient, err := helmClientFromConfig(kubeconfig, "")
	if err != nil {
		return err
	}

	slog.Info("adding Helm Chart repositories")
	repos := make(map[string]repo.Entry)
	for _, addon := range config.Spec.ClusterAddons {
		repos[addon.RepoURL] = repo.Entry{
			Name: addon.RepoName,
			URL:  addon.RepoURL,
		}
	}
	for _, r := range repos {
		slog.Info("adding Helm Chart repository", "name", r.Name, "url", r.URL)
		if err = helmClient.AddOrUpdateChartRepo(r); err != nil {
			return err
		}
	}

	// Install Helm Charts.
	for _, addon := range config.Spec.ClusterAddons {
		slog.Info("installing addon", "name", addon.Name, "namespace", addon.Namespace, "chart", addon.Chart, "version", addon.Version)

		// Reinitialize the client to match the target release namespace.
		helmClient, err = helmClientFromConfig(kubeconfig, addon.Namespace)
		if err != nil {
			return err
		}

		release := &helmclient.ChartSpec{
			ReleaseName:     addon.Name,
			ChartName:       addon.Chart,
			Namespace:       addon.Namespace,
			ValuesYaml:      addon.ValuesObject,
			CreateNamespace: true,
			UpgradeCRDs:     true,
			Wait:            true,
			Timeout:         5 * time.Minute,
		}

		_, err = helmClient.InstallOrUpgradeChart(context.Background(), release, &helmclient.GenericHelmOptions{})
		if err != nil {
			return err
		}
		slog.Info("addon installed", "name", addon.Name, "namespace", addon.Namespace, "chart", addon.Chart, "version", addon.Version)
	}
	return err
}
