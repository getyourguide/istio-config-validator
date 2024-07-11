package envoy_test

import (
	"testing"

	"github.com/getyourguide/istio-config-validator/internal/pkg/istio-router-check/envoy"
	"github.com/stretchr/testify/require"
)

func TestReadTests(t *testing.T) {
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
			parsed, err := envoy.ReadTests(tt.path)
			require.NoError(t, err)
			var gotTests []string
			for _, t := range parsed.Tests {
				gotTests = append(gotTests, t.TestName)
			}
			require.ElementsMatch(t, tt.want, gotTests)
		})
	}
}
