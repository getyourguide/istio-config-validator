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
			got, err := helpers.ReadCRDs(tt.path)
			require.NoError(t, err)
			tt.assert(t, got)
		})
	}
}

func TestReadEnvoyTests(t *testing.T) {
	for _, tt := range []struct {
		name string
		path string
		want []string
	}{{
		name: "it should return single folder test",
		path: "testdata/tests/reviews",
		want: []string{
			"test reviews.prod.svc.cluster.local/wpcatalog",
			"test reviews.prod.svc.cluster.local/consumercatalog",
		},
	}, {
		name: "it should return tests from multiple folders",
		path: "testdata/tests/",
		want: []string{
			"test details.prod.svc.cluster.local/api/v2/products",
			"test details.prod.svc.cluster.local/api/v2/items",
			"test reviews.prod.svc.cluster.local/wpcatalog",
			"test reviews.prod.svc.cluster.local/consumercatalog",
		},
	}} {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := helpers.ReadEnvoyTests(tt.path)
			require.NoError(t, err)
			var gotTests []string
			for _, t := range parsed.Tests {
				gotTests = append(gotTests, t.TestName)
			}
			require.ElementsMatch(t, tt.want, gotTests)
		})
	}
}
