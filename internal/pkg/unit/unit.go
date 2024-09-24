// Package unit contains logic to run the unit tests against istio configuration.
// It intends to replicates istio logic when it comes to matching requests and defining its destinations.
// Once the destinations are found for a given test case, it will try to assert with the expected results.
package unit

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	networking "istio.io/api/networking/v1"
	v1 "istio.io/client-go/pkg/apis/networking/v1"
)

// Run is the entrypoint to run all unit tests defined in test cases
func Run(testfiles, configfiles []string, strict bool) ([]string, []string, error) {
	var summary, details []string

	testCases, err := parser.ParseTestCases(testfiles, strict)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing testcases failed: %w", err)
	}

	virtualServices, err := parser.ParseVirtualServices(configfiles)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing virtualservices failed: %w", err)
	}

	inputCount := 0
	for _, testCase := range testCases {
		details = append(details, "running test: "+testCase.Description)
		inputs, err := testCase.Request.Unfold()
		if err != nil {
			return summary, details, err
		}
		for _, input := range inputs {
			checkHosts := true
			route, err := GetRoute(input, virtualServices, checkHosts)
			if err != nil {
				details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
				return summary, details, fmt.Errorf("error getting destinations: %v", err)
			}
			if route.Delegate != nil {
				if testCase.Delegate != nil {
					if reflect.DeepEqual(route.Delegate, testCase.Delegate) != testCase.WantMatch {
						details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
						return summary, details, fmt.Errorf("delegate missmatch=%v, want %v, rule matched: %v", route.Delegate, testCase.Delegate, route.Match)
					}
					details = append(details, fmt.Sprintf("PASS input:[%v]", input))
				}
				if testCase.Route != nil {
					vs, err := GetDelegatedVirtualService(route.Delegate, virtualServices)
					if err != nil {
						details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
						return summary, details, fmt.Errorf("error getting delegate virtual service: %v", err)
					}
					checkHosts = false
					route, err = GetRoute(input, []*v1.VirtualService{vs}, checkHosts)
					if err != nil {
						details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
						return summary, details, fmt.Errorf("error getting destinations for delegate %v: %v", route.Delegate, err)
					}
				}
			}
			if testCase.Route != nil {
				if reflect.DeepEqual(route.Route, testCase.Route) != testCase.WantMatch {
					details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
					return summary, details, fmt.Errorf("destination missmatch=%v, want %v, rule matched: %v", route.Route, testCase.Route, route.Match)
				}
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
			if testCase.Redirect != nil {
				if reflect.DeepEqual(route.Redirect, testCase.Redirect) != testCase.WantMatch {
					details = append(details, fmt.Sprintf("FAIL input:[%v]", input))
					return summary, details, fmt.Errorf("redirect missmatch=%v, want %v, rule matched: %v", route.Redirect, testCase.Redirect, route.Match)
				}
			}
			details = append(details, fmt.Sprintf("PASS input:[%v]", input))
		}
		inputCount += len(inputs)
		details = append(details, "===========================")
	}
	summary = append(summary, "Test summary:")
	summary = append(summary, fmt.Sprintf(" - %d testfiles, %d configfiles", len(testfiles), len(configfiles)))
	summary = append(summary, fmt.Sprintf(" - %d testcases with %d inputs passed", len(testCases), inputCount))
	return summary, details, nil
}

// GetRoute returns the route that matched a given input.
func GetRoute(input parser.Input, virtualServices []*v1.VirtualService, checkHosts bool) (*networking.HTTPRoute, error) {
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
					return &networking.HTTPRoute{}, err
				} else if match {
					return httpRoute, nil
				}
			}
		}
	}

	return &networking.HTTPRoute{}, nil
}

// GetDelegatedVirtualService returns the virtualservice matching namespace/name matching the delegate argument.
func GetDelegatedVirtualService(delegate *networking.Delegate, virtualServices []*v1.VirtualService) (*v1.VirtualService, error) {
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
