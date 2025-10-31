package cluster

import (
	"context"
	"log/slog"

	"github.com/datastrophic/kemu/pkg/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kwok "sigs.k8s.io/kwok/pkg/apis/v1alpha1"
)

const (
	KWOKDurationFromAnnotation = ".metadata.annotations[\"pod-complete.stage.kwok.x-k8s.io/delay\"]"
	KWOKStagePodComplete       = "pod-complete"
)

var kwokAddons = []api.ClusterAddon{
	{
		Name:      "kwok",
		RepoName:  "kwok",
		RepoURL:   "https://kwok.sigs.k8s.io/charts/",
		Namespace: "kube-system",
		Chart:     "kwok/kwok",
		Version:   "0.2.0",
	},
	{
		Name:      "kwok-stage-fast",
		RepoName:  "kwok",
		RepoURL:   "https://kwok.sigs.k8s.io/charts/",
		Namespace: "kube-system",
		Chart:     "kwok/stage-fast",
		Version:   "0.2.0",
	},
}

func InstallKWOK(kubeconfig string) error {
	slog.Info("installing kwok")
	if err := InstallOrUpgradeAddons(kwokAddons, kubeconfig); err != nil {
		return err
	}

	kwokClient, err := kwokClientFromConfig(kubeconfig)
	if err != nil {
		return err
	}

	slog.Info("updating kwok pod-complete stage")
	stage, err := kwokClient.KwokV1alpha1().Stages().Get(context.Background(), KWOKStagePodComplete, metav1.GetOptions{})
	if err != nil {
		return err
	}

	stage.Spec.Delay = &kwok.StageDelay{
		DurationFrom: &kwok.ExpressionFrom{
			JQ: &kwok.ExpressionJQ{
				Expression: KWOKDurationFromAnnotation,
			},
		},
	}

	stage, err = kwokClient.KwokV1alpha1().Stages().Update(context.Background(), stage, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	slog.Info("kwok pod-complete stage updated")
	return nil
}
