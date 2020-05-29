package parser

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ghodss/yaml"

	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

func parseVirtualServices(files []string) ([]*v1alpha3.VirtualService, error) {
	out := []*v1alpha3.VirtualService{}

	for _, file := range files {
		fileContet, err := ioutil.ReadFile(file)
		if err != nil {
			return []*v1alpha3.VirtualService{}, err
		}

		// we need to transform yaml to json so the marsheler from istio works
		jsonBytes, err := yaml.YAMLToJSON(fileContet)
		if err != nil {
			return []*v1alpha3.VirtualService{}, err
		}

		virtualService := &v1alpha3.VirtualService{}
		err = json.Unmarshal(jsonBytes, virtualService)
		if err != nil {
			return []*v1alpha3.VirtualService{}, err
		}

		if virtualService.Name == "" {
			continue
		}

		out = append(out, virtualService)
	}

	return out, nil
}
