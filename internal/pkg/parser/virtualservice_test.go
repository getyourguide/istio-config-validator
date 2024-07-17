package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

func TestParseVirtualServices(t *testing.T) {
	expectedTestCases := []*v1alpha3.VirtualService{{Spec: networkingv1alpha3.VirtualService{
		Hosts: []string{"www.example.com", "example.com"},
	}}}
	configfiles := []string{"../../../examples/virtualservice.yml"}
	virtualServices, err := ParseVirtualServices(configfiles, false)
	require.NoError(t, err)
	require.NotEmpty(t, virtualServices)

	for _, expected := range expectedTestCases {
		for _, out := range virtualServices {
			assert.ElementsMatch(t, expected.Spec.Hosts, out.Spec.Hosts)
		}
	}
}

func TestParseMultipleVirtualServices(t *testing.T) {
	expectedTestCases := []*v1alpha3.VirtualService{{Spec: networkingv1alpha3.VirtualService{
		Hosts: []string{"www.example.com", "example.com"},
	}}}

	configfiles := []string{"../../../examples/multidocument_virtualservice.yml"}
	virtualServices, err := ParseVirtualServices(configfiles, false)
	require.NoError(t, err)
	require.NotEmpty(t, virtualServices)
	require.GreaterOrEqual(t, len(virtualServices), 2)

	for _, expected := range expectedTestCases {
		for _, out := range virtualServices {
			assert.ElementsMatch(t, expected.Spec.Hosts, out.Spec.Hosts)
		}
	}
}

func TestVirtualServiceUnknownFields(t *testing.T) {
	_, err := ParseVirtualServices([]string{"testdata/invalid_vs.yml"}, true)
	require.ErrorContains(t, err, "cannot parse proto message")
}
