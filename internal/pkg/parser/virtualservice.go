package parser

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

// FIXME: Unmarshal doesn't parse the whole structure. istio/client-go seems to only parse `spec`.
func parseVirtualServices(files []string) ([]*v1alpha3.VirtualService, error) {
	out := []*v1alpha3.VirtualService{}
	for _, file := range files {
		virtualService := &v1alpha3.VirtualService{}
		fileContet, err := ioutil.ReadFile(file)
		if err != nil {
			return out, err
		}

		err = yaml.Unmarshal(fileContet, virtualService)
		if err != nil {
			return out, err
		}

		// As we don't have the whole struct filled this is the way found to check if a
		// file is a virtualservice.
		if len(virtualService.Spec.Hosts) == 0 {
			continue
		}

		out = append(out, virtualService)
	}
	return out, nil
}
