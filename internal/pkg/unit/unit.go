package unit

import (
	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	"istio.io/api/networking/v1alpha3"
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
		// Unfold input
		// for _, _ := range []parser.Input{{Authority: "www.example.com"}} {
		// }
	}
	return nil
}

func inputMatch(input parser.Input, virtualServices []*v1alpha3.VirtualService) bool {

	return true
}
