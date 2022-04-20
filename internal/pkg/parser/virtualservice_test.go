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
	testcasefiles := []string{"../../../examples/virtualservice_test.yml"}
	configfiles := []string{"../../../examples/virtualservice.yml"}
	parser, err := New(testcasefiles, configfiles)
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

func TestParseMultipleVirtualServices(t *testing.T) {
	expectedTestCases := []*v1alpha3.VirtualService{{Spec: networkingv1alpha3.VirtualService{
		Hosts: []string{"www.example.com", "example.com"},
	}}}
	testcasefiles := []string{"../../../examples/virtualservice_test.yml"}
	configfiles := []string{"../../../examples/multidocument_virtualservice.yml"}
	parser, err := New(testcasefiles, configfiles)
	if err != nil {
		t.Errorf("error getting test cases %v", err)
	}
	if len(parser.VirtualServices) == 0 {
		t.Error("virtualservices is empty")
	}
	if len(parser.VirtualServices) < 2 {
		t.Error("did not parse all virtualservices in file")
	}

	for _, expected := range expectedTestCases {
		for _, out := range parser.VirtualServices {
			assert.ElementsMatch(t, expected.Spec.Hosts, out.Spec.Hosts)
		}
	}
}
