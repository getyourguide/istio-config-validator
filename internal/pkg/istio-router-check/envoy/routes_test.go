package envoy_test

import (
	"testing"

	"github.com/getyourguide/istio-config-validator/internal/pkg/istio-router-check/envoy"
	"github.com/getyourguide/istio-config-validator/internal/pkg/istio-router-check/helpers"
	"github.com/stretchr/testify/require"
)

func TestRoutesGenerator(t *testing.T) {
	t.Run("it should generate one route", func(t *testing.T) {
		cfg, err := helpers.ReadCRDs("testdata/virtualservice.yml")
		require.NoError(t, err)
		rg := &envoy.RouteGenerator{
			Configs: cfg,
		}
		routes, err := rg.Routes()
		require.NoError(t, err)
		require.Len(t, routes, 1)
	})
}
