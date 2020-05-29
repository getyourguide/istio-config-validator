// Package parser is the package responsible for parsing test cases and istio configuration
// to be use on test assertionpackage parser
package parser

import v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"

// Parser contains the parsed files needed to run tests
type Parser struct {
	TestCases       []*TestCase
	VirtualServices []*v1alpha3.VirtualService
}

// New parses and loads the testcases and istio configuration files
func New(testfiles []string, configfiles []string) (*Parser, error) {
	testCases, err := parseTestCases(testfiles)
	if err != nil {
		return nil, err
	}

	virtualServices, err := parseVirtualServices(configfiles)
	if err != nil {
		return nil, err
	}
	parser := &Parser{
		TestCases:       testCases,
		VirtualServices: virtualServices,
	}
	return parser, nil
}
