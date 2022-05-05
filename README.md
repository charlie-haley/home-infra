# home-cluster

My personal k8s cluster built with Sidero and managed by Flux2 GitOps

## 💻 Management Cluster
| Node                     | RAM  | Storage                    | Function           | Operating System     | Quantity
| ------------------------ |------| -------------------------- | ------------------ | -------------------- | --------
| Raspberry Pi 4 Model B   | 4GB  | 256GB NVME                 | Sidero CP + Worker | Talos 1.0.3          | 1


## 💻 Metal-01 Workload Cluster
| Node                     | RAM  | Storage                    | Function           | Operating System     | Quantity
| ------------------------ |------| -------------------------- | ------------------ | -------------------- | --------
| Raspberry Pi 4 Model B   | 4GB  | 256GB NVME                 | Kube Control Plane | Talos 1.1.0-alpha.1  | 3
| Lenovo M720q             | 16GB | 256GB NVME + 1TB SSD       | Kube Worker        | Talos 1.0.3          | 3
| Raspberry Pi 4 Model B   | 4GB  | 256GB NVME                 | Kube Worker        | Talos 1.1.0-alpha.1  | 1

## Cluster

- PXE Boot Talos managed by Sidero
- Traefik Ingress
- Rook Ceph
- Prometheus/Grafana/Loki Monitoring Stack

## 🦾 Automations
- [Renovate](https://github.com/renovatebot/renovate)
- [GitHub Action YAMLlint](https://github.com/ibiqlik/action-yamllint)
- [Renovate Helm Releases](https://github.com/k8s-at-home/renovate-helm-releases)
