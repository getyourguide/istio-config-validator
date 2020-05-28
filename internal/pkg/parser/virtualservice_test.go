package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

func TestParseVirtualServices(t *testing.T) {
	expectedTestCases := []*v1alpha3.VirtualService{{Spec: networkingv1alpha3.VirtualService{
		Hosts: []string{"www.example.com", "example.com"},
	}}}
	configuration := &Configuration{
		RootDir: "../../../examples/",
	}
	parser, err := New(configuration)
	if err != nil {
		t.Errorf("error getting test cases %v", err)
	}
	if len(parser.VirtualServices) == 0 {
		t.Error("virtualservices is empty")
	}

	for _, expected := range expectedTestCases {
		for _, out := range parser.VirtualServices {
			assert.ElementsMatch(t, expected.Spec.Hosts, out.Spec.Hosts)
		}
	}
}
