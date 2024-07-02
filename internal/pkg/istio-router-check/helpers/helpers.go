package helpers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/getyourguide/istio-config-validator/internal/pkg/istio-router-check/envoy"
	"gopkg.in/yaml.v3"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"istio.io/istio/pkg/config"
)

func ReadCRDs(baseDir string) ([]config.Config, error) {
	var configs []config.Config
	yamlFiles, err := WalkYAML(baseDir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", baseDir, err)
	}
	for _, path := range yamlFiles {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", path, err)
		}
		c, _, err := crd.ParseInputs(string(data))
		if err != nil {
			return nil, fmt.Errorf("failed to parse CRD: %w\n%s", err, string(data))
		}
		configs = append(configs, c...)
	}
	return configs, nil
}

func ReadEnvoyTests(baseDir string) (envoy.Tests, error) {
	var tests envoy.Tests
	yamlFiles, err := WalkYAML(baseDir)
	if err != nil {
		return envoy.Tests{}, fmt.Errorf("error reading directory %s: %w", baseDir, err)
	}
	for _, path := range yamlFiles {
		data, err := os.ReadFile(path)
		if err != nil {
			return envoy.Tests{}, fmt.Errorf("error reading file %s: %w", path, err)
		}
		var t envoy.Tests
		err = yaml.Unmarshal(data, &t)
		if err != nil {
			return envoy.Tests{}, fmt.Errorf("error unmarshalling file %s: %w", path, err)
		}
		tests.Tests = append(tests.Tests, t.Tests...)
	}
	return tests, nil
}

// WalkYAML walks the baseDirs and returns a list of all files found with yaml extension.
func WalkYAML(baseDirs ...string) ([]string, error) {
	return WalkFilter(func(path string, info os.FileInfo) bool {
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			return true
		}
		return false
	}, baseDirs...)
}

// WalkFilter walks the baseDirs and returns a list of all files found that return true in the filterFunc.
func WalkFilter(filterFunc func(path string, info os.FileInfo) bool, baseDirs ...string) ([]string, error) {
	var files []string
	for _, baseDir := range baseDirs {
		err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error reading directory %q: %w", baseDir, err)
			}
			if info.IsDir() {
				return nil
			}
			if !filterFunc(path, info) {
				return nil
			}
			files = append(files, path)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error reading directory %q: %w", baseDir, err)
		}
	}
	return files, nil
}
