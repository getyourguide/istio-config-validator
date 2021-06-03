package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/ghodss/yaml"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
)

var (
	// ErrEmptyAuthorityList indicates an empty Authority list
	ErrEmptyAuthorityList = errors.New("authority list is empty")
	// ErrEmptyMethodList indicates an empty Method list
	ErrEmptyMethodList = errors.New("method list is empty")
	// ErrEmptyURIList indicates an empty URI list
	ErrEmptyURIList = errors.New("URI list is empty")
)

// TestCaseYAML define the list of TestCase
type TestCaseYAML struct {
	TestCases []*TestCase `yaml:"testCases"`
}

// TestCase defines the API for declaring unit tests
type TestCase struct {
	Description string                                     `yaml:"description"`
	Request     *Request                                   `yaml:"request"`
	Route       []*networkingv1alpha3.HTTPRouteDestination `yaml:"route"`
	Redirect    *networkingv1alpha3.HTTPRedirect           `yaml:"redirect"`
	Rewrite     *networkingv1alpha3.HTTPRewrite            `yaml:"rewrite"`
	Fault       *networkingv1alpha3.HTTPFaultInjection     `yaml:"fault"`
	Headers     *networkingv1alpha3.Headers                `yaml:"headers"`
	WantMatch   bool                                       `yaml:"wantMatch"`
}

// Request define the crafted http request present in the test case file.
type Request struct {
	Authority []string          `yaml:"authority"`
	Method    []string          `yaml:"method"`
	URI       []string          `yaml:"uri"`
	Headers   map[string]string `yaml:"headers"`
}

// Input contains the data structure which will be used to assert
type Input struct {
	Authority string
	Method    string
	URI       string
	Headers   map[string]string
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

// Unfold returns a list of Input objects constructed by all possibilities defined in the Request object. Ex:
// Request{Authority: {"www.example.com", "example.com"}, Method: {"GET", "OPTIONS"}}
// returns []Input{
// 	{Authority:"www.example.com", Method: "GET"},
// 	{Authority:"www.example.com", Method: "OPTIONS"}
// 	{Authority:"example.com", Method: "GET"},
// 	{Authority:"example.com", Method: "OPTIONS"},
// }
func (r *Request) Unfold() ([]Input, error) {
	out := []Input{}

	if len(r.Authority) == 0 {
		return out, ErrEmptyAuthorityList
	}
	if len(r.Method) == 0 {
		return out, ErrEmptyMethodList
	}
	if len(r.URI) == 0 {
		return out, ErrEmptyURIList
	}

	for _, uri := range r.URI {
		u, err := url.Parse(uri)
		if err != nil {
			return out, err
		}

		for _, auth := range r.Authority {
			for _, method := range r.Method {
				out = append(out, Input{Authority: auth, Method: method, URI: u.Path, Headers: r.Headers})
			}
		}
	}

	return out, nil
}

func parseTestCases(files []string) ([]*TestCase, error) {
	out := []*TestCase{}

	for _, file := range files {
		fileContent, err := ioutil.ReadFile(file)
		if err != nil {
			return []*TestCase{}, fmt.Errorf("reading file '%s' failed: %w", file, err)
		}

		// we need to transform yaml to json so the marshaler from istio works
		jsonBytes, err := yaml.YAMLToJSON(fileContent)
		if err != nil {
			return []*TestCase{}, fmt.Errorf("yamltojson conversion failed for file '%s': %w", file, err)
		}

		yamlFile := &TestCaseYAML{}
		err = json.Unmarshal(jsonBytes, yamlFile)
		if err != nil {
			return []*TestCase{}, fmt.Errorf("unmarshaling failed for file '%s': %w", file, err)
		}

		if len(yamlFile.TestCases) == 0 {
			continue
		}

		out = append(out, yamlFile.TestCases...)
	}
	return out, nil
}
