// Package unit contains logic to run the unit tests against istio configuration.
// It intends to replicates istio logic when it comes to matching requests and defining its destinations.
// Once the destinations are found for a given test case, it will try to assert with the expected results.
package unit

import (
	"fmt"
	"reflect"

	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

// Run is the entrypoint to run all unit tests defined in test cases
func Run(testfiles, configfiles []string) ([]string, []string, error) {
	var summary, details []string
	parsed, err := parser.New(testfiles, configfiles)
	if err != nil {
		return summary, details, err
	}

	inputCount := 0
	for _, testCase := range parsed.TestCases {
		details = append(details, "running test: "+testCase.Description)
		inputs, err := testCase.Request.Unfold()
		if err != nil {
			return summary, details, err
		}
		for _, input := range inputs {
			destinations, err := GetDestination(input, parsed.VirtualServices)
			if err != nil {
				details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
				return summary, details, fmt.Errorf("error getting destinations: %v", err)
			}
			if reflect.DeepEqual(destinations, testCase.Route) != testCase.WantMatch {
				details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
				return summary, details, fmt.Errorf("destination missmatch=%v, want %v", destinations, testCase.Route)
			}

			details = append(details, fmt.Sprintf("PASS input:[%v]", input))
		}
		inputCount += len(inputs)
		details = append(details, "===========================")
	}
	summary = append(summary, "Test summary:")
	summary = append(summary, fmt.Sprintf(" - %d testfiles, %d configfiles", len(testfiles), len(configfiles)))
	summary = append(summary, fmt.Sprintf(" - %d testcases with %d inputs passed", len(parsed.TestCases), inputCount))
	return summary, details, nil
}

// GetDestination return the destination list for a given input, it evaluates matching rules in VirtualServices
// related to the input based on it's hosts list.
func GetDestination(input parser.Input, virtualServices []*v1alpha3.VirtualService) ([]*networkingv1alpha3.HTTPRouteDestination, error) {
	for _, vs := range virtualServices {
		spec := vs.Spec
		if !contains(spec.Hosts, input.Authority) {
			continue
		}

		for _, httpRoute := range spec.Http {
			if len(httpRoute.Match) == 0 {
				return httpRoute.Route, nil
			}
			for _, matchBlock := range httpRoute.Match {
				if match, err := matchRequest(input, matchBlock); err != nil {
					return []*networkingv1alpha3.HTTPRouteDestination{}, err
				} else if match {
					return httpRoute.Route, nil
				}
			}
		}
	}

	return []*networkingv1alpha3.HTTPRouteDestination{}, nil
}

// GetRoute returns the route that matched a given input.
func GetRoute(input parser.Input, virtualServices []*v1alpha3.VirtualService) (*networkingv1alpha3.HTTPRoute, error) {
	for _, vs := range virtualServices {
		spec := vs.Spec
		if !contains(spec.Hosts, input.Authority) {
			continue
		}

		for _, httpRoute := range spec.Http {
			if len(httpRoute.Match) == 0 {
				return httpRoute, nil
			}
			for _, matchBlock := range httpRoute.Match {
				if match, err := matchRequest(input, matchBlock); err != nil {
					return &networkingv1alpha3.HTTPRoute{}, err
				} else if match {
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
