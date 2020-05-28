package parser

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

// FIXME: Unmarshal doesn't parse the whole structure. istio/client-go seems to only parse `spec`.
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

			virtualService := &v1alpha3.VirtualService{}
			fileContet, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			err = yaml.Unmarshal(fileContet, virtualService)
			if err != nil {
				return err
			}

			// As we don't have the whole struct filled this is the way found to check if a
			// file is a virtualservice.
			if len(virtualService.Spec.Hosts) == 0 {
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
