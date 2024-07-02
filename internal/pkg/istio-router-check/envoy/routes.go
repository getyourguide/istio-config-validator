package envoy

import (
	"fmt"

	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/test/xds"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/schema/gvk"

	istiolog "istio.io/istio/pkg/log"
	istiotest "istio.io/istio/pkg/test"
)

type RouteGenerator struct {
	Configs []config.Config
	routes  []*route.RouteConfiguration
}

// Routes returns a list of routes in Envoy format. It uses istio's fake discovery server to generate the routes.
// The routes are generated from the Configs loaded in the RouteGenerator.
// TODO(cainelli): The RouterGenerator only takes VirtualServices into account when they are bound to `mesh`, VirtualServices bound exclusively to `gateway` are not considered.
func (rg *RouteGenerator) Routes() ([]*route.RouteConfiguration, error) {
	logOpts := istiolog.DefaultOptions()
	logOpts.SetDefaultOutputLevel("all", istiolog.ErrorLevel)
	err := istiolog.Configure(logOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to configure logging: %w", err)
	}

	// Add a random endpoint, otherwise there will be no routes to check
	rg.Configs = append(rg.Configs, config.Config{
		Meta: config.Meta{
			GroupVersionKind: gvk.ServiceEntry,
			Namespace:        "a",
			Name:             "wg-a",
			Labels: map[string]string{
				"grouplabel": "notonentry",
			},
		},
		Spec: &v1alpha3.ServiceEntry{
			Hosts: []string{"pod.pod.svc.cluster.local"},
			Ports: []*v1alpha3.ServicePort{{
				Number:   80,
				Protocol: "HTTP",
				Name:     "http",
			}},
			Location:   v1alpha3.ServiceEntry_MESH_INTERNAL,
			Resolution: v1alpha3.ServiceEntry_STATIC,
			Endpoints: []*v1alpha3.WorkloadEntry{{
				Address: "10.10.10.20",
			}},
		},
	})

	err = istiotest.Wrap(func(t istiotest.Failer) {
		srv := xds.NewFakeDiscoveryServer(t, xds.FakeOptions{
			Configs: rg.Configs,
		})
		proxy := srv.SetupProxy(nil)
		rg.routes = srv.Routes(proxy)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate routes: %w", err)
	}
	return rg.routes, nil
}
