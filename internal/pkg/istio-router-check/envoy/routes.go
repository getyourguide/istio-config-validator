package envoy

import (
	"fmt"
	"os"

	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/getyourguide/istio-config-validator/internal/pkg/istio-router-check/helpers"
	"istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/test/xds"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/schema/gvk"

	istiolog "istio.io/istio/pkg/log"
	istiotest "istio.io/istio/pkg/test"
)

type routeGenerator struct {
	configs     []config.Config
	proxy       *model.Proxy
	gatewayName string
	routes      []*route.RouteConfiguration
}

func NewRouteGenerator(opts ...optionFunc) *routeGenerator {
	rg := &routeGenerator{}
	for _, opt := range opts {
		opt(rg)
	}
	return rg
}

// Routes returns a list of routes in Envoy format. It uses istio's fake discovery server to generate the routes.
// The routes are generated from the Configs loaded in the RouteGenerator.
// TODO(cainelli): The RouterGenerator only takes VirtualServices into account when they are bound to `mesh`, VirtualServices bound exclusively to `gateway` are not considered.
func (rg *routeGenerator) Routes() ([]*route.RouteConfiguration, error) {
	logOpts := istiolog.DefaultOptions()
	logOpts.SetDefaultOutputLevel("all", istiolog.ErrorLevel)
	err := istiolog.Configure(logOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to configure logging: %w", err)
	}

	// Add a random endpoint, otherwise there will be no routes to check
	rg.configs = append(rg.configs, config.Config{
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
		if err := rg.prepareProxy(); err != nil {
			t.FailNow()
		}
		srv := xds.NewFakeDiscoveryServer(t, xds.FakeOptions{
			Configs: rg.configs,
		})
		proxy := srv.SetupProxy(rg.proxy)
		rg.routes = srv.Routes(proxy)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate routes: %w", err)
	}
	return rg.routes, nil
}

// prepareProxy creates a proxy with the provided metadata. If the gateway is set, it will create a proxy with the
// metadata of the gateway. If the gateway does not exist in the provided configs, it will be created with the default
// values accepting "*" as hosts.
func (rg *routeGenerator) prepareProxy() error {
	if rg.gatewayName == "" {
		return nil
	}

	namespacedName := helpers.NewNamespacedNameFromString(rg.gatewayName)
	metadata := &model.NodeMetadata{
		Namespace: namespacedName.Namespace,
		Labels: map[string]string{
			"istio": "ingressgateway",
		},
	}
	var gatewayFound bool
	for _, cfg := range rg.configs {
		if cfg.Meta.GroupVersionKind != gvk.Gateway {
			continue
		}
		if cfg.Meta.Name == namespacedName.Name && cfg.Meta.Namespace == namespacedName.Namespace {
			gatewayFound = true
			var selector map[string]string
			switch v := cfg.Spec.(type) {
			case *v1alpha3.Gateway:
				selector = v.Selector
			default:
				return fmt.Errorf("could not cast Gateway spec (%T) for %s/%s", v, cfg.Meta.Namespace, cfg.Meta.Name)
			}

			metadata = &model.NodeMetadata{
				Namespace: cfg.Meta.Namespace,
				Labels:    selector,
			}
			break
		}
	}

	if !gatewayFound {
		rg.configs = append(rg.configs, config.Config{
			Meta: config.Meta{
				GroupVersionKind: gvk.Gateway,
				Name:             namespacedName.Name,
				Namespace:        namespacedName.Namespace,
				Labels:           metadata.Labels,
			},
			Spec: &v1alpha3.Gateway{
				Selector: metadata.Labels,
				Servers: []*v1alpha3.Server{{
					Hosts: []string{"*"},
					Port: &v1alpha3.Port{
						Number:   80,
						Protocol: "HTTP",
					},
				}},
			},
		})
	}

	rg.proxy = &model.Proxy{
		Type:     model.Router,
		Labels:   metadata.Labels,
		Metadata: metadata,
	}

	return nil
}

func ReadCRDs(baseDir string) ([]config.Config, error) {
	var configs []config.Config
	yamlFiles, err := helpers.WalkYAML(baseDir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", baseDir, err)
	}
	for _, path := range yamlFiles {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", path, err)
		}
		c, _, err := crd.ParseInputs(string(data))
		if err != nil {
			return nil, fmt.Errorf("failed to parse CRD: %w\n%s", err, string(data))
		}
		configs = append(configs, c...)
	}
	return configs, nil
}
