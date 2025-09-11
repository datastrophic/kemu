# kemu
A Kubernetes Cluster emulation tool based on Kind and KWOK.

Easily create emulated Kubernetes clusters with 1000s of nodes
for workload scheduling experimentation and analysis.

## Quickstart
Create Kind cluster:
```shell
make create-cluster
```

Install KWOK controller and Prometheus operator for observability:
```shell
make install
```

Export `KUBECONFIG`:
```shell
export KUBECONFIG=$(pwd)/.run/kubeconfig
```

Create an emulated cluster:
```shell
go run . create-cluster --cluster-config examples/gcp-cluster.yaml --kubeconfig "${KUBECONFIG}"
```

Explore created cluster (`kubectl`):
```shell
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
