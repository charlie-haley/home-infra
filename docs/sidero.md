# Sidero

The following applies to sidero v0.4
## Dependencies

- kubectl
    ```
    sudo pacman -S kubectl
    ```
- clusterctl 0.4.4
    ```
    curl -Lo /usr/local/bin/clusterctl \
    "https://github.com/kubernetes-sigs/cluster-api/releases/download/v0.4.4/clusterctl-linux-amd64"
    chmod +x /usr/local/bin/clusterctl
    ```
- talosctl
     ```
    sudo curl -Lo /usr/local/bin/talosctl \
    "https://github.com/talos-systems/talos/releases/latest/download/talosctl-linux-amd64"
    chmod +x /usr/local/bin/talosctl
     ```

## Install Talos on RPI4

Flash USB SSD with Talos image, (found here.)[https://github.com/talos-systems/talos/releases/latest/download/metal-rpi_4-arm64.img.xz]

## Creating management cluster
```bash
#ip of single-node cluster running talos
export SIDERO_ENDPOINT=192.168.1.215

#generate config
talosctl gen config --config-patch='[{"op": "add", "path": "/cluster/allowSchedulingOnMasters", "value": true},{"op": "replace", "path": "/machine/install/disk", "value": "/dev/sda"}]' rpi4-sidero https://${SIDERO_ENDPOINT}:6443/

#apply generated config
talosctl apply-config --insecure -n ${SIDERO_ENDPOINT} -f controlplane.yaml

#merge client config into ~/.talos/config
talosctl config merge talosconfig

#update config endpoints/nodes
talosctl config endpoints ${SIDERO_ENDPOINT}
talosctl config nodes ${SIDERO_ENDPOINT}

#bootstrap etcd
talosctl bootstrap

#fetch kubeconfig
talosctl kubeconfig

#wait for log "[talos] bootstrap sequence: done" before continuing
talosctl dmesg -f | grep "bootstrap sequence: done"

# seems to take a couple minutes after that log before 6443 is open and it's ready for the clusterctl command

#init management cluster
SIDERO_CONTROLLER_MANAGER_HOST_NETWORK=true SIDERO_CONTROLLER_MANAGER_API_ENDPOINT=${SIDERO_ENDPOINT} clusterctl init -i sidero -b talos -c talos

#verify admin cluster
curl -I http://${SIDERO_ENDPOINT}:8081/tftp/ipxe.efi
```

## Setting up DHCP

```bash
# example ipxe-metal.conf located in /sidero/dhcp
set service dhcp-server global-parameters 'option system-arch code 93 = unsigned integer 16;'
set service dhcp-server shared-network-name VLAN10 subnet 192.168.1.0/24 subnet-parameters "include &quot;/config/ipxe-metal.conf&quot;;"
```

## Configure servers
(follow guide to configure rpi4 as servers with PXE boot)[https://www.sidero.dev/docs/v0.4/guides/rpi4-as-servers/#build-the-image-with-the-boot-folder-contents]

## Patch metal controller
__TODO: move this into a kustomization__
(As per the documentation here)[https://www.sidero.dev/docs/v0.4/guides/rpi4-as-servers/#patch-metal-controller], we need to patch the sidero-controller-manager so the RPI4's can boot over network boot = UEFI.

```
kubectl -n sidero-system patch deployments.apps sidero-controller-manager --patch "$(cat ./sidero/cluster/core/patches/controller.patch.yaml)"
```

## Bootstrap Flux
Ensure we're using the correct context
```bash
kubectx admin@rpi4-sidero
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
kubectl apply --kustomize=./sidero/cluster/base/flux-system
```

## Get kubeconfig

```bash
# fetch talosconfig
kubectl get talosconfig -o yaml $(kubectl get talosconfig -A -o jsonpath='{.items[0].metadata.name}') -n flux-system -o jsonpath='{.status.talosConfig}' > talosconfig

# fetch kubeconfig
talosctl --nodes="192.168.1.14" --talosconfig talosconfig kubeconfig .

# merge kubeconfig files
cp ~/.kube/config ~/.kube/config.bak
KUBECONFIG=~/.kube/config:$(pwd)/kubeconfig kubectl config view --flatten > /tmp/kubeconfig
mv /tmp/kubeconfig ~/.kube/config
```

## Tidy up context names
```bash
kubectx sidero=admin@rpi4-sidero
kubectx workload=admin@metal-01
```

## Next steps

(Next we'll need to bootstrap Flux onto the workload cluster.)[metal-01.md]
