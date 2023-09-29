// Package unit contains logic to run the unit tests against istio configuration.
// It intends to replicate istio logic when it comes to matching requests and defining its destinations.
// Once the destinations are found for a given test case, it will try to assert with the expected results.
package unit

import (
	"fmt"
	"reflect"

	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
)

type Formatter interface {
	Pass(input parser.Input) string
	Fail(input parser.Input) string
}

type TestRunner struct {
	TestCases       []*parser.TestCase
	VirtualServices []*v1alpha3.VirtualService
	Format          Formatter
}

type defaultFormat struct{}

func (d defaultFormat) Fail(input parser.Input) string {
	return fmt.Sprintf("FAIL input:[%v]", input)
}
func (d defaultFormat) Pass(input parser.Input) string {
	return fmt.Sprintf("Pass input:[%v]", input)
}

var defaultFormatter Formatter = defaultFormat{}

// Run is the entrypoint to run all unit tests defined in test cases
func Run(testfiles, configfiles []string) ([]string, []string, error) {
	parsed, err := parser.New(testfiles, configfiles)
	runner := TestRunner{TestCases: parsed.TestCases, VirtualServices: parsed.VirtualServices}
	if err != nil {
		return nil, nil, err
	}
	summary, details, err := runner.Run()
	summary = append(summary, fmt.Sprintf(" - %d testfiles, %d configfiles", len(testfiles), len(configfiles)))
	return summary, details, err
}

// Run executes the TestRunner
func (tr *TestRunner) Run() ([]string, []string, error) {
	if tr.Format == nil {
		tr.Format = defaultFormatter
	}

	var summary, details []string
	inputCount := 0
	for _, testCase := range tr.TestCases {
		details = append(details, "running test: "+testCase.Description)
		inputs, err := testCase.Request.Unfold()
		if err != nil {
			return summary, details, err
		}
		for _, input := range inputs {
			route, err := GetRoute(input, tr.VirtualServices)
			if err != nil {
				details = append(details, tr.Format.Fail(input))
				return summary, details, fmt.Errorf("error getting destinations: %v", err)
			}
			if reflect.DeepEqual(route.Route, testCase.Route) != testCase.WantMatch {
				details = append(details, tr.Format.Fail(input))
				return summary, details, fmt.Errorf("destination missmatch=%v, want %v, rule matched: %v", route.Route, testCase.Route, route.Match)
			}
			if testCase.Rewrite != nil {
				if reflect.DeepEqual(route.Rewrite, testCase.Rewrite) != testCase.WantMatch {
					details = append(details, tr.Format.Fail(input))
					return summary, details, fmt.Errorf("rewrite missmatch=%v, want %v, rule matched: %v", route.Rewrite, testCase.Rewrite, route.Match)
				}
			}
			if testCase.Fault != nil {
				if reflect.DeepEqual(route.Fault, testCase.Fault) != testCase.WantMatch {
					details = append(details, tr.Format.Fail(input))
					return summary, details, fmt.Errorf("fault missmatch=%v, want %v, rule matched: %v", route.Fault, testCase.Fault, route.Match)
				}
			}
			if testCase.Headers != nil {
				if reflect.DeepEqual(route.Headers, testCase.Headers) != testCase.WantMatch {
					details = append(details, tr.Format.Fail(input))
					return summary, details, fmt.Errorf("headers missmatch=%v, want %v, rule matched: %v", route.Headers, testCase.Headers, route.Match)
				}
			}

			details = append(details, tr.Format.Pass(input))
		}
		inputCount += len(inputs)
		details = append(details, "===========================")
	}
	summary = append(summary, "Test summary:")
	summary = append(summary, fmt.Sprintf(" - %d testcases with %d inputs passed", len(tr.TestCases), inputCount))
	return summary, details, nil
}

// GetRoute returns the route that matched a given input.
func GetRoute(input parser.Input, virtualServices []*v1alpha3.VirtualService) (*networkingv1alpha3.HTTPRoute, error) {
	for _, vs := range virtualServices {
		spec := &vs.Spec
		if !contains(spec.Hosts, input.Authority) {
			continue
		}

		for _, httpRoute := range spec.Http {
			if len(httpRoute.Match) == 0 {
				return httpRoute, nil
			}
			for _, matchBlock := range httpRoute.Match {
				match, err := matchRequest(input, matchBlock)
				if err != nil {
					return &networkingv1alpha3.HTTPRoute{}, err
				}
				if match {
					return httpRoute, nil
				}
			}
		}
	}

	return &networkingv1alpha3.HTTPRoute{}, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
