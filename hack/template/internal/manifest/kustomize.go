package manifest

import (
	"os"

	kustomize "sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/yaml"
)

func (m *Manifest) ProcessKustomize(app string, namespace string, appDir string) error {
	k := &kustomize.Kustomization{
		Namespace: namespace,
		Resources: m.Kustomize,
	}
	resource, err := yaml.Marshal(&k)
	if err != nil {
		return err
	}
	err = os.WriteFile(appDir+"/kustomization.yaml", resource, 0755)
	if err != nil {
		return err
	}

	return nil
}
