# CNI

This folder contains the manifests needed to deploy the CNI (Cilium). Due to v1.9+ of Cilium removing the single URL containing deployment manifests, this has to be templated out and hosted on GitHub. https://www.talos.dev/v1.0/kubernetes-guides/network/deploying-cilium/

This is referenced as patches in the machineconfig for Talos.

```yaml
- op: replace
  path: /cluster/network/cni
  value:
    name: "custom"
    urls:
      - "https://raw.githubusercontent.com/charlie-haley/home-cluster/main/integrations/sidero/custom/cni/cilium-deployment.yaml
```
