#!/bin/bash
set -e

KUSTOMIZATIONS_CHART="./templates/framework/"
MANIFESTS_DIR="./manifests/"
tmpdir=`mktemp -d`

ks_values=""

if [[ ! -z "${PUBLISH}" ]]; then
  printf "🟧 Detected publish mode...\n"
  printf "🟧 Logging into container registry...\n"
  echo $PASS | docker login $REGISTRY -u $USER --password-stdin
  printf "🟩 Logged into container registry.\n"
fi

create_kustomize () {
  cat <<EOF  > $app_dir/kustomization.yaml
namespace: $namespace
resources:
$1
EOF
}

# Loop over each directory in the manifests directory
for d in $MANIFESTS_DIR/*/ ; do
    namespace=`basename $d`
    printf "🟧 Processing namespace $namespace...\n"

    ns_dir="$tmpdir/$namespace/"
    mkdir $ns_dir

    ks_ns_val=`cat << EOM
  $namespace:
EOM`

    # For each application
    for f in $d*; do
      values=`cat $f | yq .values`
      helm=`cat $f | yq .helm`
      kustomize=`cat $f | yq -r -o="yaml" .kustomize`
      resources=`cat $f | yq -r -o="yaml" ".resources[]"`

      repo=`cat $f | yq .helm.repo | tr -d \"`
      chart=`cat $f | yq .helm.chart`
      version=`cat $f | yq .helm.version`

      # This is a bit gross... need to migrate to an actual programming language soon™
      dependsOn=`cat $f | yq -r -o="yaml" .dependsOn | sed  "s/^/      /"`
      fileName=`basename $f`
      release=${fileName%%.*}

      printf "🟧 Processing app $release...\n"

      app_dir="$tmpdir/$namespace/$release"
      mkdir $app_dir

      if [[ "$helm" != "null" ]]; then
        helm repo add $release $repo >/dev/null
        helm repo update >/dev/null
        template=`helm template $release $release/$chart --no-hooks --version $version --include-crds --kube-version="1.27" -a "monitoring.coreos.com/v1","networking.k8s.io/v1" --values -  <<EOF
$values
EOF`
        echo "$template" >> "$app_dir/chart.yaml"
        create_kustomize "- chart.yaml"
      fi

      if [[ "$kustomize" != "null" ]]; then
        # If kustomization exists, append resources
        kustomize_file="$app_dir/kustomization.yaml"
        if test -f $kustomize_file; then
          printf "$kustomize" >> $kustomize_file
        fi

        # Kustomization doesn't exist, create it and append resources
        create_kustomize "$kustomize"
      fi

      if [[ ! -z "$resources" ]]; then
        res_kustomize_files=""
        readarray resourceArray < <(cat $f | yq e -o=j -I=0 '.resources[]')
        for rs in "${resourceArray[@]}"; do
          res_name=`printf "$rs" | yq .metadata.name`
          res_kind=`printf "$rs" | yq .kind`
          resource_filename="$res_name-$res_kind.yaml"

          printf "$rs" | yq -o="yaml" -P | yq '.. style="double"' | grep "" > "$app_dir/$resource_filename"
          res_kustomize_files="$res_kustomize_files\n- $resource_filename"
        done

        res_kustomize_files=`printf "$res_kustomize_files"`
        # If kustomization exists, append resources
        kustomize_file="$app_dir/kustomization.yaml"
        if test -f $kustomize_file; then
          printf "$res_kustomize_files" >> $kustomize_file
        fi
        # Kustomization doesn't exist, create it and append resources
        create_kustomize "$res_kustomize_files"
      fi

      # Apppend variable to the end of the framework chart values
      # If it has dependencies, make sure they're specified
      if [[ "$dependsOn" =~ .*"null".* ]]; then
        ks_ns_val="$ks_ns_val\n    $release: {}"
      else
        ks_ns_val="$ks_ns_val\n    $release:\n      dependsOn:\n$dependsOn"
      fi

      printf "🟩 Processed app $release\n"
    done

    ks_values="$ks_values\n$ks_ns_val"
    printf "🟩 Processed namespace $namespace\n"
done

# If publish mode, push artifacts to registry
if [[ ! -z "${PUBLISH}" ]]; then
  tag="${GIT_SHA:0:7}-$(date +%s)"

  # Push Flux manifests artifact
  flux push artifact oci://$IMAGE/manifests:$tag \
    --source="https://github.com/charlie-haley/home-infra" \
    --revision=$GIT_SHA \
    --path $tmpdir/

  input=`printf "namespaces:$ks_values" | yq`

  ks_artifact=`helm template kustomizations $KUSTOMIZATIONS_CHART --values - <<EOF
$input
EOF`

  # Push Flux kustomizations artifact
  flux push artifact oci://$IMAGE/kustomizations:$tag \
    --source="https://github.com/charlie-haley/home-infra" \
    --revision=$GIT_SHA \
    --path - <<EOF
$ks_artifact
EOF
else
  # Print values passed to framework chart
  input=`printf "namespaces:$ks_values" | yq`
  echo "$input"
fi
