// Package parser is the package responsible for parsing test cases and istio configuration
// to be use on test assertion
package parser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Configuration needed to build test cases
type Configuration struct {
	RootDir string
}

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

// New returns the list of test cases for a given configuration
func New(configuration *Configuration) ([]*TestCase, error) {
	testCases, err := parseTestCases(configuration.RootDir)
	if err != nil {
		return nil, err
	}
	return testCases, nil
}

func parseTestCases(rootDir string) ([]*TestCase, error) {
	out := []*TestCase{}
	err := filepath.Walk(rootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.HasSuffix(path, "_test.yml") || strings.HasSuffix(path, "_test.yaml") {
				yamlFile := &TestCaseYAML{}
				fileContet, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				err = yaml.Unmarshal(fileContet, yamlFile)
				if err != nil {
					return err
				}

				out = append(out, yamlFile.TestCases...)
			}

			return nil
		})

	if err != nil {
		return nil, err
	}
	return out, nil
}
