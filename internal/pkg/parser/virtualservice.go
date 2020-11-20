package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"go.uber.org/zap/zapcore"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"istio.io/pkg/log"
)

func parseVirtualServices(files []string) ([]*v1alpha3.VirtualService, error) {
	out := []*v1alpha3.VirtualService{}

	for _, file := range files {
		fileContent, err := ioutil.ReadFile(file)
		if err != nil {
			return []*v1alpha3.VirtualService{}, fmt.Errorf("reading file '%s' failed: %w", file, err)
		}

		// we need to transform yaml to json so the marshaler from istio works
		jsonBytes, err := yaml.YAMLToJSON(fileContent)
		if err != nil {
			log.Debug("error converting yaml to json", zapcore.Field{Key: "file", Type: zapcore.StringType, String: file})
			continue
		}

		meta := &v1.TypeMeta{}
		err = json.Unmarshal(jsonBytes, meta)
		if err != nil {
			log.Debug("error extracting the metadata of the virtualservice", zapcore.Field{Key: "file", Type: zapcore.StringType, String: file})
			continue
		}

		if meta.Kind != "VirtualService" {
			log.Debug("file is not Kind VirtualService", zapcore.Field{Key: "file", Type: zapcore.StringType, String: file})
			continue
		}

		virtualService := &v1alpha3.VirtualService{}
		err = json.Unmarshal(jsonBytes, virtualService)
		if err != nil {
			log.Debug("error while trying to unmarshal virtualservice", zapcore.Field{Key: "file", Type: zapcore.StringType, String: file})
			return out, fmt.Errorf("error while trying to unmarshal virtualservice (%s): %w", file, err)
		}

		out = append(out, virtualService)
	}

	return out, nil
}
