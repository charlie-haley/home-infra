# home-infra

<p align="center" style="text-align: center">
    <img src="./docs/images/gopher.svg" width="30%"><br/>
<br/>
    <strong>One Repo to Rule Them All.</strong> <br/>
  <i>A monorepo for managing my home infrastructure using GitOps.</i><br/>
</p>

## 📁 Folder Structure

```
├── docs           # Documentation as markdown files.
├── hack           # Scripts and other bits
├── kubernetes     # Kubernets manifests
│   ├── bootstrap    # Manifests required when bootstrapping the cluster for the first time
│   ├── manifests    # Application deployment manifests
│   └── templates    # Local Helm templates
```

## ☸️ Kubernetes

I run a bare metal cluster provisioned using [Talos Linux](https://www.talos.dev/) and managed using [Flux](https://fluxcd.io/). The cluster is comprised of 3 worker nodes and 1 control plane node.

| Hostname  | Node              | Resources          |
| --------- | ----------------- | ------------------ |
| vilya-c01 | Lenovo M720q Tiny | 16GB RAM, i5-9500t |
| vilya-w01 | Lenovo M720q Tiny | 16GB RAM, i5-9500t |
| vilya-w02 | Lenovo M720q Tiny | 16GB RAM, i5-9500t |
| vilya-w03 | Lenovo M720q Tiny | 16GB RAM, i5-9500t |

### ⚙️ GitOps

I use Flux to manage deployments to the cluster, everything that is deployed to my cluster is defined as YAML files in the `kubernetes/manifests/` directory.

```
├── manifests                # Manifests deployed to the cluster
│   ├── cert-manager      # The namespace for all the files in the directory to be deployed to
│   ├── home
│   └── storage
└── gitops              # Anything and everything Flux/GitOps related
    ├── flux-system
```

To save on duplicate code and reduce the management overhead of adding new Flux Kustomization's, Namespaces and other boilerplate config. I template all of my manifests and deploy them up to GHCR as an OCI image. This is done through a (pretty hacky) bash script which reads the provided YAML file and actions it based on the contents.

The script is triggered by Github Actions, a webhook is then fired after the package is uploaded which tells Flux to reconcile the cluster with the state from the OCI image.

## 👏 Thanks

Thanks to everyone in the [Kubernetes@Home Discord community](https://discord.gg/k8s-at-home) Inspiration for how to deploy and my manage my cluster has been influenced heavily by everyone who's shared their clusters using the [k8s-at-home GitHub tag](https://github.com/topics/k8s-at-home).
