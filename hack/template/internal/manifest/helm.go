package manifest

import (
    "time"

    helmv2 "github.com/fluxcd/helm-controller/api/v2beta2"
    sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime/schema"
)

func (m *Manifest) ProcessHelm(app string, namespace string, appDir string) error {
    // Create HelmRelease
    hr := &helmv2.HelmRelease{
        ObjectMeta: metav1.ObjectMeta{
            Name:      app,
            Namespace: namespace,
        },
        Spec: helmv2.HelmReleaseSpec{
            Interval: metav1.Duration{
                Duration: 15 * time.Minute,
            },
            Chart: helmv2.HelmChartTemplate{
                Spec: helmv2.HelmChartTemplateSpec{
                    Chart:   m.Helm.Chart,
                    Version: m.Helm.Version,
                    SourceRef: helmv2.CrossNamespaceObjectReference{
                        Kind: sourcev1.HelmRepositoryKind,
                        Name: app,
                    },
                },
            },
            Values: m.Values,
        },
    }
    hr.SetGroupVersionKind(schema.GroupVersionKind{Group: "helm.toolkit.fluxcd.io", Version: "v2beta2", Kind: "HelmRelease"})

    err := createFile(hr, "/helm-release.yaml", appDir)
    if err != nil {
        return err
    }
    m.Kustomize = append(m.Kustomize, "helm-release.yaml")

    // Create HelmRepository
    repo := &sourcev1.HelmRepository{
        ObjectMeta: metav1.ObjectMeta{
            Name:      app,
            Namespace: namespace,
        },
        Spec: sourcev1.HelmRepositorySpec{
            Interval: metav1.Duration{
                Duration: 15 * time.Minute,
            },
            URL: m.Helm.Repo,
        },
    }
    repo.SetGroupVersionKind(schema.GroupVersionKind{Group: "source.toolkit.fluxcd.io", Version: "v1beta2", Kind: "HelmRepository"})

    err = createFile(repo, "/helm-repository.yaml", appDir)
    if err != nil {
        return err
    }
    m.Kustomize = append(m.Kustomize, "helm-repository.yaml")

    return nil
}
