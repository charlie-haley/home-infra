package manifest

import (
    "os"

    fluxmetav1 "github.com/fluxcd/pkg/apis/meta"
    apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/cli-runtime/pkg/printers"
)

type Manifest struct {
    Helm        *Helm                                  `json:"helm,omitempty"`
    DependsOn   []fluxmetav1.NamespacedObjectReference `json:"dependsOn"`
    Values      *apiextensionsv1.JSON                  `json:"values,omitempty"`
    ValuesFrom  *apiextensionsv1.JSON                  `json:"valuesFrom,omitempty"`
    Resources   []*apiextensionsv1.JSON                `json:"resources,omitempty"`
    Kustomize   []string                               `json:"kustomize,omitempty"`
    Backup      *Backup                                `json:"backup,omitempty"`
}

type Helm struct {
    Repo    string `json:"repo"`
    Chart   string `json:"chart"`
    Version string `json:"version"`
}

type Backup struct {
    Pvc string `json:"pvc"`
}

func (m *Manifest) Process(app string, namespace string, appDir string) error {
    if m.Values != nil {
        err := m.ProcessHelm(app, namespace, appDir)
        if err != nil {
            return err
        }
    }

    if m.Resources != nil {
        err := m.ProcessResources(app, namespace, appDir)
        if err != nil {
            return err
        }
    }

    if m.Backup != nil {
        err := m.ProcessBackup(app, namespace, appDir)
        if err != nil {
            return err
        }
    }

    return m.ProcessKustomize(app, namespace, appDir)
}

func createFile(obj runtime.Object, filename string, appDir string) error {
    file, err := os.Create(appDir + filename)
    if err != nil {
        return err
    }

    defer file.Close()

    y := printers.YAMLPrinter{}
    return y.PrintObj(obj, file)
}
