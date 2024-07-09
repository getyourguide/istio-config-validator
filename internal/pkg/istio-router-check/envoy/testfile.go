package envoy

import (
	"fmt"
	"os"

	"github.com/getyourguide/istio-config-validator/internal/pkg/istio-router-check/helpers"
	"gopkg.in/yaml.v3"
)

type Tests struct {
	Tests []Test `yaml:"tests,omitempty" json:"tests,omitempty"`
}
type Test struct {
	TestName string   `yaml:"test_name,omitempty" json:"test_name,omitempty"`
	Input    Input    `yaml:"input,omitempty" json:"input,omitempty"`
	Validate Validate `yaml:"validate" json:"validate"`
}

type Validate struct {
	ClusterName           string                `yaml:"cluster_name" json:"cluster_name"`
	VirtualClusterName    string                `yaml:"virtual_cluster_name,omitempty" json:"virtual_cluster_name,omitempty"`
	VirtualHostName       string                `yaml:"virtual_host_name,omitempty" json:"virtual_host_name,omitempty"`
	HostRewrite           string                `yaml:"host_rewrite,omitempty" json:"host_rewrite,omitempty"`
	PathRewrite           string                `yaml:"path_rewrite,omitempty" json:"path_rewrite,omitempty"`
	PathRedirect          string                `yaml:"path_redirect,omitempty" json:"path_redirect,omitempty"`
	RequestHeaderMatches  []RequestHeaderMatch  `yaml:"request_header_matches,omitempty" json:"request_header_matches,omitempty"`
	ResponseHeaderMatches []ResponseHeaderMatch `yaml:"response_header_matches,omitempty" json:"response_header_matches,omitempty"`
}

type RequestHeaderMatch struct {
	Name        string      `yaml:"name,omitempty" json:"name,omitempty"`
	StringMatch StringMatch `yaml:"string_match,omitempty" json:"string_match,omitempty"`
}

type ResponseHeaderMatch struct {
	Name          string         `yaml:"name,omitempty" json:"name,omitempty"`
	StringMatch   map[string]any `yaml:"string_match,omitempty" json:"string_match,omitempty"`
	PresenceMatch PresenceMatch  `yaml:"presence_match,omitempty" json:"presence_match,omitempty"`
}

type StringMatch struct {
	Exact string `yaml:"exact,omitempty" json:"exact,omitempty"`
}

type PresenceMatch struct{}

type Input struct {
	Authority                 string   `yaml:"authority,omitempty" json:"authority,omitempty"`
	Path                      string   `yaml:"path,omitempty" json:"path,omitempty"`
	Method                    string   `yaml:"method,omitempty" json:"method,omitempty"`
	Internal                  bool     `yaml:"internal,omitempty" json:"internal,omitempty"`
	RandomValue               string   `yaml:"random_value,omitempty" json:"random_value,omitempty"`
	SSL                       bool     `yaml:"ssl,omitempty" json:"ssl,omitempty"`
	Runtime                   string   `yaml:"runtime,omitempty" json:"runtime,omitempty"`
	AdditionalRequestHeaders  []Header `yaml:"additional_request_headers,omitempty" json:"additional_request_headers,omitempty"`
	AdditionalResponseHeaders []Header `yaml:"additional_response_headers,omitempty" json:"additional_response_headers,omitempty"`
}

type Header struct {
	Key   string `yaml:"key,omitempty" json:"key,omitempty"`
	Value string `yaml:"value,omitempty" json:"value,omitempty"`
}

func ReadTests(baseDir string) (Tests, error) {
	var tests Tests
	yamlFiles, err := helpers.WalkYAML(baseDir)
	if err != nil {
		return Tests{}, fmt.Errorf("error reading directory %s: %w", baseDir, err)
	}
	for _, path := range yamlFiles {
		data, err := os.ReadFile(path)
		if err != nil {
			return Tests{}, fmt.Errorf("error reading file %s: %w", path, err)
		}
		var t Tests
		err = yaml.Unmarshal(data, &t)
		if err != nil {
			return Tests{}, fmt.Errorf("error unmarshalling file %s: %w", path, err)
		}
		tests.Tests = append(tests.Tests, t.Tests...)
	}
	return tests, nil
}
