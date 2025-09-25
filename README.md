# kemu
A Kubernetes cluster Emulator based on [Kind (Kubernetes IN Docker)](https://kind.sigs.k8s.io/)
and [KWOK (Kubernetes WithOut Kubelet)](https://kwok.sigs.k8s.io/).

Easily create large-scale emulated Kubernetes clusters with 1000s of nodes
for workload scheduling experimentation and analysis with minimal hardware
requirements.

The goal of the project is to provide a fully automated reproducible bootstrap
of emulated clusters based on declarative specification.

Follow the [Quickstart](#quickstart) guide to see KEMU in action, and the
[Overview](#kemu-overview) section for additional details and configuration walkthrough.

## Quickstart
#### Prerequisites
* [Go](https://go.dev/doc/install)
* [Docker](https://docs.docker.com/engine/install)
* [Kind](https://kind.sigs.k8s.io/)
* [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
* [Helm](https://helm.sh/docs/intro/install/)

> NOTE: Following steps should be run from the root of this project.

#### Install `kemu`
```shell
go install ./...
```

#### Create a cluster
```shell
kemu create-cluster --cluster-config examples/gcp-cluster.yaml --kubeconfig $(pwd)/kemu.config
```

#### Explore the cluster
```shell
export KUBECONFIG=$(pwd)/kemu.config

kubectl get nodes

# Example output:
# NAME                    STATUS   ROLES           AGE     VERSION
# a2-ultragpu-8g-use1-0   Ready    agent           99s     v1.33.1
# a2-ultragpu-8g-use1-1   Ready    agent           99s     v1.33.1
# a2-ultragpu-8g-use1-2   Ready    agent           99s     v1.33.1
# a2-ultragpu-8g-use1-3   Ready    agent           99s     v1.33.1
# a2-ultragpu-8g-use1-4   Ready    agent           99s     v1.33.1
# a2-ultragpu-8g-use2-0   Ready    agent           99s     v1.33.1
# a2-ultragpu-8g-use2-1   Ready    agent           99s     v1.33.1
# a2-ultragpu-8g-use2-2   Ready    agent           99s     v1.33.1
# a2-ultragpu-8g-use2-3   Ready    agent           99s     v1.33.1
# a2-ultragpu-8g-use2-4   Ready    agent           99s     v1.33.1
# a2-ultragpu-8g-use3-0   Ready    agent           99s     v1.33.1
# a2-ultragpu-8g-use3-1   Ready    agent           99s     v1.33.1
# a2-ultragpu-8g-use3-2   Ready    agent           99s     v1.33.1
# a2-ultragpu-8g-use3-3   Ready    agent           99s     v1.33.1
# a2-ultragpu-8g-use3-4   Ready    agent           98s     v1.33.1
# a3-highgpu-8g-use1-0    Ready    agent           98s     v1.33.1
# a3-highgpu-8g-use1-1    Ready    agent           98s     v1.33.1
# a3-highgpu-8g-use1-2    Ready    agent           98s     v1.33.1
# a3-highgpu-8g-use1-3    Ready    agent           98s     v1.33.1
# a3-highgpu-8g-use1-4    Ready    agent           97s     v1.33.1
# a3-highgpu-8g-use2-0    Ready    agent           97s     v1.33.1
# a3-highgpu-8g-use2-1    Ready    agent           97s     v1.33.1
# a3-highgpu-8g-use2-2    Ready    agent           97s     v1.33.1
# a3-highgpu-8g-use2-3    Ready    agent           97s     v1.33.1
# a3-highgpu-8g-use2-4    Ready    agent           96s     v1.33.1
# a3-ultragpu-8g-use1-0   Ready    agent           96s     v1.33.1
# a3-ultragpu-8g-use1-1   Ready    agent           96s     v1.33.1
# a3-ultragpu-8g-use1-2   Ready    agent           96s     v1.33.1
# a3-ultragpu-8g-use1-3   Ready    agent           96s     v1.33.1
# a3-ultragpu-8g-use1-4   Ready    agent           95s     v1.33.1
# kwok-control-plane      Ready    control-plane   7m58s   v1.33.1
```

#### Delete the cluster
```shell
kemu delete-cluster
```

## KEMU Overview
KEMU provides a single-spec declarative approach for configuring control plane nodes,
installing cluster addons, and defining emulated cluster nodes with various capacity
and placement options.

KEMU builds on Kind for control plane and worker nodes deployment which are used for
running auxiliary software required for the experimentation. Examples include Prometheus
Operator for observability, custom schedulers (Volcano, Yunikorn), workload management
operators (Kueue, KubeRay), etc. Running these components requires actual Kubelet(s) to
be available for scheduling, and Kind provides sufficient functionality for this.

To provide a reproducible cluster dependencies setup and configuration, KEMU supports
cluster addons defined as Helm Charts. A `ClusterConfig` spec allows specifying a list
of Helm Charts that are installed on cluster bootstrap automatically. Each cluster
addon can be parametrized via `valuesObject` with the same content as the standard
Helm values file.

Kubelet emulation is based on KWOK, and KEMU provides a lightweight configuration
scheme for defining node groups with various properties, and generates specified
number of nodes automatically.

The [Example Configuration](#example-configuration) section provides a `ClusterConfig`
specification walkthrough and explanation of core configuration properties. 

### Example Configuration
The following example demonstrates key components of the KEMU `ClusterConfig` specification.
The specification contains 3 main sections:
* `kindConfig` - a YAML configuration file used for creating Kind Cluster. This is a standard
  [Kind Configuration](https://kind.sigs.k8s.io/docs/user/configuration/) which is passed to
  Kind cluster provisioner without any modifications.
* `clusterAddons` define a list of Helm Charts to be installed as a part of cluster
  bootstrap process. Each cluster addon can be provided with `valuesObject` containing
  Helm Chart values for the installation.
* `nodeGroups` define groups of emulated nodes sharing similar properties (instance type, capacity)
  and the placement of the nodes. Node placement allows configuring number of nodes in different
  availability zones.

Example specification:
```yaml
apiVersion: kemu.datastrophic.io/v1alpha1
kind: ClusterConfig
spec:
  kindConfig: |
    kind: Cluster
    apiVersion: kind.x-k8s.io/v1alpha4
    nodes:
      - role: control-plane
      - role: worker
      - role: worker
  clusterAddons:
    - name: kwok
      repoName: kwok
      repoURL: https://kwok.sigs.k8s.io/charts/
      namespace: kube-system
      chart: kwok/kwok
      version: 0.2.0
    - name: kwok-stage-fast
      repoName: kwok
      repoURL: https://kwok.sigs.k8s.io/charts/
      namespace: kube-system
      chart: kwok/stage-fast
      version: 0.2.0
    - name: prometheus
      repoName: prometheus-community
      repoURL: https://prometheus-community.github.io/helm-charts
      namespace: monitoring
      chart: prometheus-community/kube-prometheus-stack
      version: 75.16.1
      valuesObject: |
        alertmanager:
          enabled: false
  nodeGroups:
    - name: a2-ultragpu-8g
      placement:
        - availabilityZone: use1
          replicas: 5
        - availabilityZone: use2
          replicas: 5
        - availabilityZone: use3
          replicas: 5
      nodeTemplate:
        metadata:
          labels:
            datastrophic.io/gpu-type: nvidia-a100-80gb
        capacity:
          cpu: 96
          memory: 1360Gi
          ephemeralStorage: 3Ti
          nvidia.com/gpu: 8
```

## Development
KEMU relies on end-to-end and integration tests to verify its functionality.
Tests require Kind, Docker, and Helm being installed on the machine where they run.
A KEMU cluster is created based on provided configuration by each of the tests.

Run parallel tests with Ginkgo:
```shell
go install github.com/onsi/ginkgo/v2/ginkgo@latest

ginkgo -v -p ./...
```

Run tests with coverage:
```shell
ginkgo -v -p -coverprofile=coverage.out -coverpkg=./pkg/... ./...
go tool cover -html=coverage.out -o coverage.html
```
