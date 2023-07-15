# 🪙 Sidero

## Dependencies

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
- [helm](https://helm.sh/docs/helm/helm_install/)
- [talosctl](https://www.talos.dev/v1.4/learn-more/talosctl/)

## ⚙️ Generate Talos Config

Before we provision the nodes for the cluster, we will need to generate the Talos config. We can do this using existing secrets that are encrypted in the repo. Before we can use the secrets they will need decrypting with SOPS. If you don't have existing secrets or want to create new ones, this can be done with `talosctl gen secrets -o kubernetes/bootstrap/talos/secrets.yaml`

```bash
sops --in-place --decrypt kubernetes/bootstrap/talos/secrets.yaml
```

Once they're decrypted, we're able to generate the config using `talosctl`. This also includes patches to modify the default configuration for both control plane and worker node types. This will generate all the config into a temporary directory with the suffix `-TALOS`. **it's important to clear this up once done as it contains secrets**

```bash
export CONTROL_PLANE_HOST="vilya-c01"
export TALOSCONFIG_DIR=$(mktemp -d --suffix=-TALOS)

talosctl gen config \
    --with-secrets kubernetes/talos/secrets.yaml \
    vilya \
    https://$CONTROL_PLANE_HOST:6443 \
    --config-patch-control-plane @kubernetes/bootstrap/talos/patch/control-plane.yaml \
    --config-patch-worker @kubernetes/bootstrap/talos/patch/worker.yaml \
    --output $TALOSCONFIG_DIR
```

We can now set the generated Talos config to be our default local config by using the `config merge` command. We will also need to set the host of the node in the Talosconfig, this can be done easily with the `config endpoint` command.

```bash
# set talosconfig as local default
talosctl config merge $TALOSCONFIG_DIR/talosconfig

# add control plane host to config
talosctl config endpoint $CONTROL_PLANE_HOST
```

## ✅Applying the Talos Config

Now we've generated all the config we can apply it to our respective nodes, starting with the control plane. We'll need to flash a bootable USB with the latest Talos image before we can start applying config. [The image can be found under the latest GitHub releases](https://github.com/siderolabs/talos/releases/latest), the image name is `talos-amd64.iso`.

After plugging the USB into the control plane node and booting from the USB, we can applying the Talos configuration to the node.

```bash
talosctl apply-config --insecure \
    --nodes $CONTROL_PLANE_HOST \
    --file $TALOSCONFIG_DIR/controlplane.yaml
```

We can bootstrap the node once the configuration is applied, this will form the Kubernetes cluster. After the server has finished bootstrapping we can fetch the Kubeconfig, by default, this will merge with the existing local configuration file.

```bash
# bootstrap cluster
talosctl bootstrap --nodes $CONTROL_PLANE_HOST

# fetch kubeconfig
talosctl kubeconfig --nodes $CONTROL_PLANE_HOST

# update context name if using kubectx
kubectx vilya=admin@vilya
```

## 🕸Applying the CNI

[As per the recommendation from Talos, when deploying Cilium as a CNI, you should deploy Talos without a CNI.](https://www.talos.dev/v1.4/kubernetes-guides/network/deploying-cilium/). The Control Plane node will be stuck on pending and the node will be stuck in a 10 minute reboot loop. During this reboot loop, we can apply some basic Cilium configuration to get the cluster online. Once online, and Flux is applied, we will start managing the CNI through Flux.

If not already added, we'll need to add the Helm repository for Cilium.

```bash
helm repo add cilium https://helm.cilium.io
helm repo update
```

Once added, we can install the Cilium chart using the values defined in the repo. As we're installing without `kube-proxy` it's important to use to use the provided values file.

```bash
helm install -n kube-system cilium cilium/cilium -f kubernetes/talos/cni/cilium-values.yaml
```

## 💼 Adding Worker Nodes

Now that the cluster is ready, the CNI has been installed and we have all our Talos config. We can add worker nodes to the cluster using the config we generated earlier. This will need doing for each worker that we'd like adding to the cluster. _(Note: multiple nodes can be passed into the --nodes arg to make it easier)_

```bash
talosctl apply-config --insecure \
    --nodes vilya-w01 \
    --file $TALOSCONFIG_DIR/worker.yaml
```

## 👞 Bootstrap Flux

Once the cluster is in a ready state we can bootstrap Flux so we can beging managing the deployment of manifests in this repo. Instead of using `flux bootstrap` the installation will be done through `kubectl` due to the repo already containing the Flux configuration.

```bash
kubectl apply -f https://github.com/fluxcd/flux2/releases/download/v2.0.1/install.yaml
```

Once Flux is running we can apply the secret needed for SOPS decryption, then, apply the configuration in the repo.

```bash
# add age key for SOPS
kubectl create secret generic sops-age \
    --namespace=flux-system \
    --from-file=age.agekey=flux.agekey

# apply Flux configuration from repo
kubectl apply -k kubernetes/manifests/gitops/
```

## ✨ Tidy Up

Now we've provisioned everything, it's important to tidy up our dangling config and re-encrypt our secrets so we don't accidentally expose them to the world!

```bash
# encrypt talos secrets
sops --in-place --encrypt kubernetes/bootstrap/talos/secrets.yaml

# tidy config
rm -R $TALOSCONFIG_DIR
```

It's also worth configuring pre-commit hooks to prevent accidentally commiting decrypted secrets. It also provides some other linting features.

```bash
sudo pacman -S python-pre-commit
pre-commit install
```

The cluster should now be provisioned 🎉
