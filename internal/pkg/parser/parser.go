// Package parser is the package responsible for parsing test cases and istio configuration
// to be use on test assertionpackage parser
package parser

import v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"

// Configuration needed to build test cases
type Configuration struct {
	RootDir string
}

// Parsed contains the parsed files needed to run tests
type Parsed struct {
	TestCases       []*TestCase
	VirtualServices []*v1alpha3.VirtualService
}

// New returns the list of test cases for a given configuration
func New(configuration *Configuration) (*Parsed, error) {
	testCases, err := parseTestCases(configuration.RootDir)
	if err != nil {
		return nil, err
	}

	virtualServices, err := parseVirtualServices(configuration.RootDir)
	if err != nil {
		return nil, err
	}
	parser := &Parsed{
		TestCases:       testCases,
		VirtualServices: virtualServices,
	}
	return parser, nil
}
