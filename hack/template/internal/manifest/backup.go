package manifest

import (
	volsyncv1 "github.com/backube/volsync/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/api/resource"
)

func (m *Manifest) ProcessBackup(app string, namespace string, appDir string) error {
	a := map[string]string{"kustomize.toolkit.fluxcd.io/prune": "disabled"}

	resticRepo := app + "-restic-config"
	s := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        resticRepo,
			Namespace:   namespace,
			Annotations: a,
		},
		StringData: map[string]string{
			"RESTIC_REPOSITORY":     "${RESTIC_REPOSITORY}/" + namespace + "/" + app,
			"RESTIC_PASSWORD":       "${RESTIC_PASSWORD}",
			"AWS_ACCESS_KEY_ID":     "${AWS_ACCESS_KEY_RESTIC}",
			"AWS_SECRET_ACCESS_KEY": "${AWS_SECRET_KEY_RESTIC}",
			"AWS_DEFAULT_REGION":    "eu-west-1",
		},
	}
	s.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"})
	err := createFile(s, "/restic-config-secret.yaml", appDir)
	if err != nil {
		return err
	}
	m.Kustomize = append(m.Kustomize, "restic-config-secret.yaml")

	schedule := "0 3 * * *"
	retainDays := int32(14)
	rs := &volsyncv1.ReplicationSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:        app,
			Namespace:   namespace,
			Annotations: a,
		},
		Spec: volsyncv1.ReplicationSourceSpec{
			SourcePVC: m.Backup.Pvc,
			Trigger: &volsyncv1.ReplicationSourceTriggerSpec{
				Schedule: &schedule,
			},
			Restic: &volsyncv1.ReplicationSourceResticSpec{
				Unlock:            "false",
				PruneIntervalDays: &retainDays,
				Repository:        resticRepo,
				CacheCapacity:     resource.MustParse("5Gi"),
				Retain: &volsyncv1.ResticRetainPolicy{
					Daily: &retainDays,
				},
				ReplicationSourceVolumeOptions: volsyncv1.ReplicationSourceVolumeOptions{
					CopyMethod: volsyncv1.CopyMethodSnapshot,
				},
			},
		},
	}
	rs.SetGroupVersionKind(schema.GroupVersionKind{Group: "volsync.backube", Version: "v1alpha1", Kind: "ReplicationSource"})

	err = createFile(rs, "/replication-source.yaml", appDir)
	if err != nil {
		return err
	}
	m.Kustomize = append(m.Kustomize, "replication-source.yaml")

	return nil
}
