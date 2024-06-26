package helpers_test

import (
	"testing"

	"github.com/getyourguide/istio-config-validator/internal/pkg/istio-router-check/helpers"
	"github.com/stretchr/testify/require"
	"istio.io/istio/pkg/config"
)

func TestReadCRDs(t *testing.T) {
	for _, tt := range []struct {
		name   string
		path   string
		assert func(t *testing.T, got []config.Config)
	}{{
		name: "it should return single virtualservice",
		path: "testdata/config-files/reviews",
		assert: func(t *testing.T, got []config.Config) {
			require.Len(t, got, 1)
			require.Equal(t, "reviews-route", got[0].GetName())
		},
	}, {
		name: "it should return multiple virtualservices across directories",
		path: "testdata/config-files",
		assert: func(t *testing.T, got []config.Config) {
			require.Len(t, got, 3)
			wantVS := []string{"reviews-route", "product-details-route", "details-fallback"}
			for _, c := range got {
				require.Contains(t, wantVS, c.GetName())
			}
		},
	}} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := helpers.ReadCRDs(tt.path)
			require.NoError(t, err)
			tt.assert(t, got)
		})
	}
}

func TestReadEnvoyTests(t *testing.T) {
	for _, tt := range []struct {
		name   string
		path   string
		assert func(t *testing.T, got helpers.EnvoyTests)
	}{{
		name: "it should return single test",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := helpers.ReadTests(tt.path)
			require.NoError(t, err)
			tt.assert(t, got)
		})
	}
}
