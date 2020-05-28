package parser

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"

	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

func parseVirtualServices(rootDir string) ([]*v1alpha3.VirtualService, error) {
	out := []*v1alpha3.VirtualService{}
	err := filepath.Walk(rootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			fileContet, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// we need to transform yaml to json so the marsheler from istio works
			jsonBytes, err := yaml.YAMLToJSON(fileContet)
			if err != nil {
				return err
			}

			virtualService := &v1alpha3.VirtualService{}
			err = json.Unmarshal(jsonBytes, virtualService)
			if err != nil {
				return err
			}

			if virtualService.Name == "" {
				return nil
			}

			out = append(out, virtualService)

			return nil
		})
	if err != nil {
		return nil, err
	}
	return out, nil
}
