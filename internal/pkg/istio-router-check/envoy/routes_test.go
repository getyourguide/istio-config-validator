package envoy_test

import (
	"testing"

	"github.com/getyourguide/istio-config-validator/internal/pkg/istio-router-check/envoy"
	"github.com/stretchr/testify/require"
	"istio.io/istio/pkg/config"
)

func TestRoutesGenerator(t *testing.T) {
	t.Run("it should generate one route", func(t *testing.T) {
		cfg, err := envoy.ReadCRDs("testdata/virtualservice.yml")
		require.NoError(t, err)
		rg := envoy.NewRouteGenerator(
			envoy.WithConfigs(cfg),
			envoy.WithGateway("istio-system/istio-ingressgateway"),
		)
		routes, err := rg.Routes()
		require.NoError(t, err)
		require.Len(t, routes, 1)
	})
}

func TestReadCRDs(t *testing.T) {
	for _, tt := range []struct {
		name   string
		path   string
		assert func(t *testing.T, got []config.Config)
	}{{
		name: "it should return single virtualservice",
		path: "testdata/bookinfo/reviews",
		assert: func(t *testing.T, got []config.Config) {
			require.Len(t, got, 1)
			require.Equal(t, "reviews-route", got[0].GetName())
		},
	}, {
		name: "it should return multiple virtualservices across directories",
		path: "testdata/bookinfo",
		assert: func(t *testing.T, got []config.Config) {
			require.Len(t, got, 3)
			wantVS := []string{"reviews-route", "product-details-route", "details-fallback"}
			for _, c := range got {
				require.Contains(t, wantVS, c.GetName())
			}
		},
	}} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := envoy.ReadCRDs(tt.path)
			require.NoError(t, err)
			tt.assert(t, got)
		})
	}
}
