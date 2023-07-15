# 🪙 Sidero

The following documentation applies to Talos v1.4 and Sidero

## Dependencies

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
- [clusterctl](https://cluster-api.sigs.k8s.io/user/quick-start.html#install-clusterctl)
- [talosctl](https://www.talos.dev/v0.9/introduction/getting-started/#talosctl)

## Provisioning the Management Cluster

Flash a USB with the latest Talos image found under [Talos GitHub Release's](https://github.com/siderolabs/talos/releases).

Next, once the device has booted and is running the Talos image we can configure the management cluster.

```bash
export SIDERO_ENDPOINT=192.168.1.215

# decrypt the secrets bundle
sops --decrypt --in-place hack/talos/config/secrets.yaml

# generate config for the management cluster
talosctl gen config \
    --config-patch @hack/talos/config/patch-management.yaml \
    --with-secrets hack/talos/config/secrets.yaml \
    management \
    https://$SIDERO_ENDPOINT:6443

# re-encrypt secrets bundle after generating config
sops --encrypt --in-place hack/talos/config/secrets.yaml

# apply generated config
talosctl apply-config --insecure -n 192.168.1.215 --file controlplane.yaml

# move talos config
mv talosconfig ~/.talos/config

# bootstrap Talos
talosctl bootstrap -n 192.168.1.215

# get kubeconfig
talosctl kubeconfig -n 191.168.1.215

# tidy
rm controlplane.yaml worker.yaml
```

## Install Sidero

Now the management cluster is running, we can install Sidero using `clusterctl`.

```
export SIDERO_CONTROLLER_MANAGER_HOST_NETWORK=true
export SIDERO_CONTROLLER_MANAGER_API_ENDPOINT=$SIDERO_ENDPOINT
export SIDERO_CONTROLLER_MANAGER_SIDEROLINK_ENDPOINT=$SIDERO_ENDPOINT

clusterctl init -b talos -c talos -i sidero
```

## Setting up DHCP

```bash
# example ipxe-metal.conf located in ./integrations/sidero/dhcp
set service dhcp-server global-parameters 'option system-arch code 93 = unsigned integer 16;'
set service dhcp-server shared-network-name VLAN10 subnet 192.168.1.0/24 subnet-parameters "include &quot;/config/ipxe-metal.conf&quot;;"
```




## Bootstrap Flux
Ensure we're using the correct context
```bash
kubectx admin@management
```
Run pre-installation checks
```bash
flux check --pre
```
Create flux-system namespace
```bash
kubectl create namespace flux-system
```
Add Age key for SOPS
```bash
cat flux.agekey | kubectl create secret generic sops-age \
    --namespace=flux-system \
    --from-file=age.agekey=/dev/stdin
```
Install Flux
```bash
# due to a race condition with the Flux CRDs, this command will need to be run twice
kubectl apply --kustomize=./manifests/management/gitops/flux-system
```

## Get kubeconfig

```bash
# fetch kubeconfig
kubectl get secret -n flux-system metal-01-kubeconfig -o yaml -o jsonpath='{.data.value}' | base64 -d > kubeconfig

# merge kubeconfig files
cp ~/.kube/config ~/.kube/config.bak
KUBECONFIG=~/.kube/config:$(pwd)/kubeconfig kubectl config view --flatten > /tmp/kubeconfig
mv /tmp/kubeconfig ~/.kube/config
```

## Tidy up context names
```bash
kubectx management=admin@management
kubectx metal-01=metal-01-admin@metal-01
```
