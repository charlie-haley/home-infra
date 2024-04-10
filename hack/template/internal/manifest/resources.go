package manifest

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"sigs.k8s.io/yaml"
)

type RemoteResource struct {
	Url    *string `json:"url"`
	Sha256 *string `json:"256sum"`
}

func (m *Manifest) ProcessResources(app string, namespace string, appDir string) error {
	for i, k := range m.Resources {
		resource, err := yaml.Marshal(&k)
		if err != nil {
			return err
		}

		// If the resource type is remote, process it accordingly
		var remote RemoteResource
		yaml.Unmarshal(resource, &remote)
		if remote.Url != nil && remote.Sha256 != nil {
			resp, err := http.Get(string(*remote.Url))
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, resp.Body)
			if err != nil {
				return err
			}

			// Validate checksum of remote resource
			sum := sha256.Sum256(buf.Bytes())
			if hex.EncodeToString(sum[:]) != *remote.Sha256 {
				return errors.New("ERROR: checksum of remote file doesn't match")
			}

			resource = buf.Bytes()
		}

		fileName := fmt.Sprintf("%v-", i) + app + ".yaml"
		err = os.WriteFile(appDir+"/"+fileName, resource, 0755)
		if err != nil {
			return err
		}

		m.Kustomize = append(m.Kustomize, fileName)
	}
	return nil
}
