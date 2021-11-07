# home-cluster

My personal k8s cluster built with Sidero and managed by Flux2 GitOps

## 💻 Nodes
| Node                     | RAM  | Storage       | Function           | Operating System
| ------------------------ |------| ------------- | ------------------ | -------------------- |
| Raspberry Pi 4 Model B   | 4GB  | 256GB NVME    | Sidero Master Node | Talos 0.13.2         |
| Raspberry Pi 4 Model B   | 4GB  | 256GB NVME    | Kube Master Node   | Talos 0.13.2         |
| Raspberry Pi 4 Model B   | 4GB  | 120GB SSD     | Kube Master Node   | Talos 0.13.2         |
| Raspberry Pi 4 Model B   | 4GB  | 120GB SSD     | Kube Master Node   | Talos 0.13.2         |
| Dell R210II              | 16GB | 1TB SSD       | Kube Worker Node   | Talos 0.13.2         |
| Dell R210II              | 16GB | 1TB SSD       | Kube Worker Node   | Talos 0.13.2         |
| Raspberry Pi 4 Model B   | 4GB  | 120GB SSD     | Kube Worker Node   | Talos 0.13.2         |
| HP MicroServer G8        | 8GB  | x4 3TB WD Red | NFS Server         | Ubuntu 20.04.2 LTS   |

## Cluster

- PXE Boot Talos managed by Sidero
- Traefik Ingress
- Mayastor
- Prometheus/Grafana/Loki Monitoring Stack

## 🦾 Automations
- [Renovate](https://github.com/renovatebot/renovate)
- [GitHub Action YAMLlint](https://github.com/ibiqlik/action-yamllint)
- [Renovate Helm Releases](https://github.com/k8s-at-home/renovate-helm-releases)
