SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec
.DEFAULT_GOAL := help

CONFIG_DIR = $(CURDIR)/config
SCRIPT_DIR = $(CURDIR)/hack
RUN_DIR = $(CURDIR)/.run
GENERATED_DIR = $(RUN_DIR)/generated
KUBECONFIG = $(RUN_DIR)/kubeconfig

# Software versions
HELM_VERSION = v3.17.4
KIND_VERSION = v0.29.0
KWOK_CHART_VERSION = 0.2.0 # https://artifacthub.io/packages/helm/kwok/kwok
PROMETHEUS_CHART_VERSION = 75.16.1 # https://github.com/prometheus-community/helm-charts/releases

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Tools Installationkind exam
LOCALBIN ?= $(RUN_DIR)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

# Binaries
KUBECTL ?= kubectl
KIND ?= $(LOCALBIN)/kind
HELM ?= $(LOCALBIN)/helm

.PHONY: kind
kind: $(KIND) ## Download Kind CLI locally if not present
$(KIND): $(LOCALBIN)
	test -s $(LOCALBIN)/kind || \
	curl -Lo $(LOCALBIN)/kind https://kind.sigs.k8s.io/dl/$(KIND_VERSION)/kind-$(shell go env GOOS)-$(shell go env GOARCH) && \
	chmod +x $(KIND)

.PHONY: helm
helm: $(HELM)  ## Download Helm CLI locally if not present
$(HELM): $(LOCALBIN)
	test -s $(LOCALBIN)/helm || \
	curl https://get.helm.sh/helm-$(HELM_VERSION)-$(shell go env GOOS)-$(shell go env GOARCH).tar.gz | tar -xz --strip-components 1 -C $(LOCALBIN) && \
	chmod +x $(HELM)

##@ Cluster Bootstrap

.PHONY: create-cluster
create-cluster: kind  ## Create Kind cluster
	$(KIND) create cluster --name kwok --config $(CONFIG_DIR)/cluster.yaml --kubeconfig=$(KUBECONFIG)

.PHONY: delete-cluster
delete-cluster: kind  ## Delete Kind cluster
	$(KIND) delete cluster --name kwok

##@ Cluster Dependencies

.PHONY: prometheus
prometheus: helm  ## Install Prometheus stack
	$(HELM) repo add prometheus-community https://prometheus-community.github.io/helm-charts
	$(HELM) repo update
	$(HELM) --kubeconfig=$(KUBECONFIG) upgrade --namespace monitoring --create-namespace --install \
		prometheus prometheus-community/kube-prometheus-stack --version $(PROMETHEUS_CHART_VERSION) -f $(CONFIG_DIR)/prometheus-values.yaml

.PHONY: kwok
kwok: helm  ## Install KWOK control plane
	$(HELM) repo add kwok https://kwok.sigs.k8s.io/charts/
	$(HELM) repo update
	$(HELM) --kubeconfig=$(KUBECONFIG) upgrade --namespace kube-system --install \
		kwok kwok/kwok --version $(KWOK_CHART_VERSION)
	$(HELM) --kubeconfig=$(KUBECONFIG) upgrade --namespace kube-system --install \
    		kwok-stage-fast kwok/stage-fast --version $(KWOK_CHART_VERSION)
	$(HELM) --kubeconfig=$(KUBECONFIG) upgrade --namespace kube-system --install \
        	kwok-metrics-usage kwok/metrics-usage --version $(KWOK_CHART_VERSION)

.PHONY: install
install: prometheus kwok  ## Install all dependencies on the cluster

##@ Helper Targets
.PHONY: grafana
grafana: ## Run kubectl port-forward to Grafana Pod
	$(KUBECTL) --kubeconfig=$(KUBECONFIG) port-forward --namespace monitoring svc/prometheus-grafana 8080:80
