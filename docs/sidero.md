# Sidero

The following applies to sidero v0.4
## Dependencies

- kubectl
    ```
    sudo pacman -S kubectl
    ```
- clusterctl 0.4.4
    ```
    curl -L https://github.com/kubernetes-sigs/cluster-api/releases/download/v0.4.4/clusterctl-linux-amd64 -o clusterctl
    chmod +x ./clusterctl
    sudo mv ./clusterctl /usr/local/bin/clusterctl
    ```
- talosctl
     ```
    sudo curl -Lo /usr/local/bin/talosctl \
    "https://github.com/talos-systems/talos/releases/latest/download/talosctl-linux-amd64"
    chmod +x /usr/local/bin/talosctl
     ```
    
## Creating management cluster
```bash
#ip of single-node cluster running talos
export SIDERO_ENDPOINT=192.168.1.215

#generate config
talosctl gen config --config-patch='[{"op": "add", "path": "/cluster/allowSchedulingOnMasters", "value": true},{"op": "replace", "path": "/machine/install/disk", "value": "/dev/sda0"}]' rpi4-sidero https://${SIDERO_ENDPOINT}:6443/

#apply generated config
talosctl apply-config --insecure -n ${SIDERO_ENDPOINT} -f controlplane.yml

#merge client config into ~/.talos/config
talosctl config merge talosconfig

#update config endpoints/nodes
talosctl config endpoints ${SIDERO_ENDPOINT}
talosctl config nodes ${SIDERO_ENDPOINT}

#bootstrap etcd
talosctl bootstrap

#fetch kubeconfig
talosctl kubeconfig

#init management cluster
SIDERO_CONTROLLER_MANAGER_HOST_NETWORK=true SIDERO_CONTROLLER_MANAGER_API_ENDPOINT=${SIDERO_ENDPOINT} clusterctl init -i sidero -b talos -c talos

#verify admin cluster
curl -I http://${SIDERO_ENDPOINT}:8081/tftp/ipxe.efi
```

## Setting up DHCP

```bash
# example ipxe-metal.conf located in /sidero/dhcp
set service dhcp-server shared-network-name VLAN10 subnet 192.168.1.0/24 subnet-parameters "include &quot;/config/ipxe-metal.conf&quot;;"
```

## Configure servers
patch server to set install disk to /dev/sda
```bash
cd sidero/patches/
kubectl patch server <server-id> --type merge --patch "$(cat server.patch.yaml)"
```

(follow guide to configure rpi4 as servers with PXE boot)[https://www.sidero.dev/docs/v0.3/guides/rpi4-as-servers/#build-the-image-with-the-boot-folder-contents]

extra step needed for RPI4 2 step boot procedure, install raspbian and run to update the EEPROM
```bash
sudo -E rpi-eeprom-config --edit
# add entry to bottom of the file - TFTP_IP: 192.168.1.215
sudo reboot
```

## Configure environments
```bash
cd sidero/environments/
kubectl apply -f default.yaml
kubectl apply -f default-arm664.yaml
```
### Accept servers
```bash
# get servers
kubectl get servers -o wide

# accept servers
kubectl get servers <server-id> --type='json' -p='[{"op": "replace", "path": "/spec/accepted", "value": true}]'
```

```
export CONTROL_PLANE_SERVERCLASS=any
export WORKER_SERVERCLASS=any
export TALOS_VERSION=v0.11
export KUBERNETES_VERSION=v1.21.3
export CONTROL_PLANE_PORT=6443
export CONTROL_PLANE_ENDPOINT=192.168.1.225
clusterctl config cluster management-plane -i sidero > management-plane.yaml
```



