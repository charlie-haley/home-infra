# Workload Cluster

The following applies to sidero v0.4
## Dependencies

- kubectl
    ```
    sudo pacman -S kubectl
    ```

## Reqirements

(Ensure the cluster has been created following the documentation here.)[sidero.md]

## Bootstrap Flux
Ensure we're using the correct context
```bash
kubectx workload
```
Run pre-installation checks
```bash
flux check --pre
```
Create flux-system namespace
```bash
kubectl create namespace flux-system
```
Add GPG key for SOPS
```bash
export FLUX_FINGERPRINT=9BED42A6B950B27737E31539730EBA837FB2813F
gpg --export-secret-keys --armor "${FLUX_FINGERPRINT}" |
kubectl create secret generic sops-gpg \
    --namespace=flux-system \
    --from-file=sops.asc=/dev/stdin
```
Install Flux
```bash
# due to a race condition with the Flux CRDs, this command will need to be run twice
kubectl apply --kustomize=./cluster/base/flux-system
```
