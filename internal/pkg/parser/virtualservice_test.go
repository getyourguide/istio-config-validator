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
	virtualServices, err := ParseVirtualServices(configfiles)
	require.NoError(t, err)
	require.NotEmpty(t, virtualServices)

	for _, expected := range expectedTestCases {
		for _, out := range virtualServices {
			assert.ElementsMatch(t, expected.Spec.Hosts, out.Spec.Hosts)
		}
	}
}

func TestParseMultipleVirtualServices(t *testing.T) {
	wantHosts := []string{"www.example2.com", "example2.com", "www.example3.com", "example3.com"}

	configfiles := []string{"../../../examples/multidocument_virtualservice.yml"}
	virtualServices, err := ParseVirtualServices(configfiles)
	require.NoError(t, err)
	require.NotEmpty(t, virtualServices)
	require.GreaterOrEqual(t, len(virtualServices), 2)

	var gotHosts []string
	for _, vs := range virtualServices {
		gotHosts = append(gotHosts, vs.Spec.Hosts...)
	}
	require.ElementsMatch(t, wantHosts, gotHosts)
}

func TestVirtualServiceUnknownFields(t *testing.T) {
	vsFiles := []string{"testdata/invalid_vs.yml"}
	_, err := ParseVirtualServices(vsFiles)
	require.ErrorContains(t, err, "cannot parse proto message")
}
