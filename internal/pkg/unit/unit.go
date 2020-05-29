package unit

import (
	"fmt"
	"reflect"

	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"istio.io/pkg/log"
)

// Configuration needed for unit tests
type Configuration struct {
	RootDir string
}

// Run is the entrypoint to run all unit tests defined in test cases
func Run(configuration *Configuration) error {
	c := &parser.Configuration{
		RootDir: configuration.RootDir,
	}
	parsed, err := parser.New(c)
	if err != nil {
		return err
	}

	for _, testCase := range parsed.TestCases {
		log.Infof("running test: %s", testCase.Description)
		inputs, err := testCase.Request.Unfold()
		if err != nil {
			return err
		}
		for _, input := range inputs {
			destinations := GetDestination(input, parsed.VirtualServices)
			if reflect.DeepEqual(destinations, testCase.Route) != testCase.WantMatch {
				return fmt.Errorf("Destination missmatch=%v, want %v", destinations, testCase.Route)
			}

			log.Infof("input:[%v] PASS", input)
		}
	}
	return nil
}

// GetDestination return the destination list for a given input, it evaluates matching rules in VirtualServices
// related to the input based on it's hosts list.
func GetDestination(input parser.Input, virtualServices []*v1alpha3.VirtualService) []*networkingv1alpha3.HTTPRouteDestination {
	for _, vs := range virtualServices {
		spec := vs.Spec
		if !contains(spec.Hosts, input.Authority) {
			continue
		}

		for _, httpRoute := range spec.Http {
			for _, matchBlock := range httpRoute.Match {
				if matchRequest(input, matchBlock) {
					return httpRoute.Route
				}
			}
		}
	}

	return []*networkingv1alpha3.HTTPRouteDestination{}
}

// AssertDestination will return true if the expected destination is present on the destination from matching rule
func AssertDestination(out []*networkingv1alpha3.HTTPRouteDestination, expected *parser.Destination) bool {
	return false
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
