package parser

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// TestCaseYAML define the list of TestCase
type TestCaseYAML struct {
	TestCases []*TestCase `yaml:"testCases"`
}

// TestCase defines the API for declaring unit tests
type TestCase struct {
	Description string       `yaml:"description"`
	Request     *Request     `yaml:"request"`
	Destination *Destination `yaml:"destination"`
}

// Request define the crafted http request to be used as input for the test case
type Request struct {
	Authority []string          `yaml:"authority"`
	Method    []string          `yaml:"method"`
	URI       []string          `yaml:"uri"`
	Headers   map[string]string `yaml:"headers"`
}

// Destination define the destination we should assert
type Destination struct {
	Host string `yaml:"host"`
	Port Port   `yaml:"port"`
}

// Port define the port of a given Destination
type Port struct {
	Number int16 `yaml:"number"`
}

func parseTestCases(rootDir string) ([]*TestCase, error) {
	out := []*TestCase{}
	err := filepath.Walk(rootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			yamlFile := &TestCaseYAML{}
			fileContet, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			err = yaml.Unmarshal(fileContet, yamlFile)
			if err != nil {
				return err
			}

			if len(yamlFile.TestCases) == 0 {
				return nil
			}

			out = append(out, yamlFile.TestCases...)

			return nil
		})

	if err != nil {
		return nil, err
	}
	return out, nil
}
