package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
	yamlV3 "gopkg.in/yaml.v3"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ParseVirtualServices(files []string) ([]*v1alpha3.VirtualService, error) {
	out := []*v1alpha3.VirtualService{}

	for _, file := range files {
		fileContent, err := os.ReadFile(file)
		if err != nil {
			return []*v1alpha3.VirtualService{}, fmt.Errorf("reading file '%s' failed: %w", file, err)
		}

		decoder := yamlV3.NewDecoder(strings.NewReader(string(fileContent)))

		for {
			// Reading into interface first. Decoding directly into struct does not work for Uri StringMatch types
			var vsInterface interface{}

			if err = decoder.Decode(&vsInterface); err != nil {
				// We've read every document in the file and we can break out
				if err == io.EOF {
					break
				}

				slog.Debug("error while trying to unmarshal into interface", zapcore.Field{Key: "file", Type: zapcore.StringType, String: file})
				return out, fmt.Errorf("error while trying to unmarshal into interface (%s): %w", file, err)
			}

			jsonBytes, err := json.Marshal(vsInterface)
			if err != nil {
				slog.Debug("error while trying to marshal to json", zapcore.Field{Key: "file", Type: zapcore.StringType, String: file})
				return out, fmt.Errorf("error while trying to marshal to json (%s): %w", file, err)
			}

			meta := &v1.TypeMeta{}
			if err = json.Unmarshal(jsonBytes, meta); err != nil {
				slog.Debug("error extracting the metadata of the virtualservice", zapcore.Field{Key: "file", Type: zapcore.StringType, String: file})
				continue
			}

			if meta.Kind != "VirtualService" {
				slog.Debug("file is not Kind VirtualService", zapcore.Field{Key: "file", Type: zapcore.StringType, String: file})
				continue
			}

			virtualService := &v1alpha3.VirtualService{}
			if err = json.Unmarshal(jsonBytes, virtualService); err != nil {
				slog.Debug("error while trying to unmarshal virtualservice", zapcore.Field{Key: "file", Type: zapcore.StringType, String: file})
				return out, fmt.Errorf("error while trying to unmarshal virtualservice (%s): %w", file, err)
			}

			out = append(out, virtualService)
		}
	}

	return out, nil
}
