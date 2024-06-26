package helpers

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"istio.io/istio/pkg/config"
)

func ReadCRDs(baseDir string) ([]config.Config, error) {
	var configs []config.Config
	err := filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error reading directory %s: %w", baseDir, err)
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}
		c, _, err := crd.ParseInputs(string(data))
		if err != nil {
			return fmt.Errorf("failed to parse CRD: %w\n%s", err, string(data))
		}

		configs = append(configs, c...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", baseDir, err)
	}
	return configs, nil
}

// TODO(cainelli): Write a tool to migrate tests from our format to Envoy's format.
type EnvoyTests struct {
	Tests []EnvoyTest `yaml:"tests" json:"tests"`
}

type EnvoyTest struct {
	TestName string         `yaml:"test_name" json:"test_name"`
	Input    map[string]any `yaml:"input" json:"input"`
	Validate map[string]any `yaml:"validate" json:"validate"`
}

func ReadTests(baseDir string) (EnvoyTests, error) {
	var tests EnvoyTests
	err := filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error reading directory %s: %w", baseDir, err)
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}
		var t EnvoyTests
		err = yaml.Unmarshal(data, &t)
		if err != nil {
			return fmt.Errorf("error unmarshalling file %s: %w", path, err)
		}
		tests.Tests = append(tests.Tests, t.Tests...)
		return nil
	})
	if err != nil {
		return EnvoyTests{}, fmt.Errorf("error reading directory %s: %w", baseDir, err)
	}
	return tests, nil
}
