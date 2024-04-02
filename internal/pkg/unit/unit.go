// Package unit contains logic to run the unit tests against istio configuration.
// It intends to replicates istio logic when it comes to matching requests and defining its destinations.
// Once the destinations are found for a given test case, it will try to assert with the expected results.
package unit

import (
	"fmt"
	"reflect"
	"slices"

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
			checkHosts := true
			route, err := GetRoute(input, parsed.VirtualServices, checkHosts)
			if err != nil {
				details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
				return summary, details, fmt.Errorf("error getting destinations: %v", err)
			}
			if route.Delegate != nil {
				// If delegate is defined in the test case, assume the delegated virtual service is not in the parsed VirtualServices list.
				// Compare the delegate configuration and skip the rest of the checks.
				if testCase.Delegate != nil {
					if reflect.DeepEqual(route.Delegate, testCase.Delegate) != testCase.WantMatch {
						details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
						return summary, details, fmt.Errorf("delegate missmatch=%v, want %v, rule matched: %v", route.Delegate, testCase.Delegate, route.Match)
					}
					details = append(details, fmt.Sprintf("PASS input:[%v]", input))
					continue
				}

				// Lookup delegated virtual service and run the rest of the checks as usual.
				vs, err := GetDelegatedVirtualService(route.Delegate, parsed.VirtualServices)
				if err != nil {
					details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
					return summary, details, fmt.Errorf("error getting delegate virtual service: %v", err)
				}
				checkHosts = false
				route, err = GetRoute(input, []*v1alpha3.VirtualService{vs}, checkHosts)
				if err != nil {
					details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
					return summary, details, fmt.Errorf("error getting destinations: %v", err)
				}
			}
			if reflect.DeepEqual(route.Route, testCase.Route) != testCase.WantMatch {
				details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
				return summary, details, fmt.Errorf("destination missmatch=%v, want %v, rule matched: %v", route.Route, testCase.Route, route.Match)
			}
			if testCase.Rewrite != nil {
				if reflect.DeepEqual(route.Rewrite, testCase.Rewrite) != testCase.WantMatch {
					details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
					return summary, details, fmt.Errorf("rewrite missmatch=%v, want %v, rule matched: %v", route.Rewrite, testCase.Rewrite, route.Match)
				}
			}
			if testCase.Fault != nil {
				if reflect.DeepEqual(route.Fault, testCase.Fault) != testCase.WantMatch {
					details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
					return summary, details, fmt.Errorf("fault missmatch=%v, want %v, rule matched: %v", route.Fault, testCase.Fault, route.Match)
				}
			}
			if testCase.Headers != nil {
				if reflect.DeepEqual(route.Headers, testCase.Headers) != testCase.WantMatch {
					details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
					return summary, details, fmt.Errorf("headers missmatch=%v, want %v, rule matched: %v", route.Headers, testCase.Headers, route.Match)
				}
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

// GetRoute returns the route that matched a given input.
func GetRoute(input parser.Input, virtualServices []*v1alpha3.VirtualService, checkHosts bool) (*networkingv1alpha3.HTTPRoute, error) {
	for _, vs := range virtualServices {
		spec := &vs.Spec
		if checkHosts && !slices.Contains(spec.Hosts, input.Authority) {
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

func GetDelegatedVirtualService(delegate *networkingv1alpha3.Delegate, virtualServices []*v1alpha3.VirtualService) (*v1alpha3.VirtualService, error) {
	for _, vs := range virtualServices {
		if vs.Name == delegate.Name {
			if delegate.Namespace != "" && vs.Namespace != delegate.Namespace {
				continue
			}
			return vs, nil
		}
	}
	return nil, fmt.Errorf("virtualservice %s not found", delegate.Name)
}
