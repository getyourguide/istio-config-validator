package envoy

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/golang/protobuf/jsonpb"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/xds"
	v3 "istio.io/istio/pilot/pkg/xds/v3"
	"istio.io/istio/pilot/test/xdstest"
	"istio.io/istio/pkg/config"
	"istio.io/pkg/log"
)

func Generate(proxyType model.NodeType, configs []config.Config) {
	s, proxy := SetupTest("", configs)

	// To determine which routes to generate, first gen listeners once (not part of benchmark) and extract routes
	l := s.Discovery.ConfigGenerator.BuildListeners(proxy, s.PushContext())
	routeNames := xdstest.ExtractRoutesFromListeners(l)
	if len(routeNames) == 0 {
		log.Fatal("Got no route names!")
	}
	generator, ok := s.Discovery.Generators[v3.RouteType]
	if !ok {
		log.Fatal("cannot find generator for %s", v3.RouteType)
	}
	c := generator.Generate(proxy, s.PushContext(), &model.WatchedResource{ResourceNames: routeNames}, nil)
	for i, r := range c {
		s, err := (&jsonpb.Marshaler{}).MarshalToString(r)
		if err != nil {
			log.Fatal(err)
		}
		// workaround to remove part of string generated that is not accepted by route config check tool
		// find a better way for doing it.
		s = strings.Replace(s, `"@type":"type.googleapis.com/envoy.config.route.v3.RouteConfiguration",`, "", 1)

		fileName := fmt.Sprintf("/tmp/examples/route-%v.json", routeNames[i])
		err = ioutil.WriteFile(fileName, []byte(s), 0644)
		if err != nil {
			log.Fatal(err)
		}

		// Cannot use b.Logf, it truncates
		fmt.Printf("Generated: %s", fileName)

	}

}

// SetupTest test builds a mock test environment. Note: push context is not initialized, to be able to benchmark separately
// most should just call setupAndInitializeTest
func SetupTest(proxyType model.NodeType, configs []config.Config) (*xds.FakeDiscoveryServer, *model.Proxy) {
	if proxyType == "" {
		proxyType = model.SidecarProxy
	}
	proxy := &model.Proxy{
		Type:        proxyType,
		IPAddresses: []string{"1.1.1.1"},
		ID:          "v0.default",
		DNSDomain:   "default.example.org",
		Metadata: &model.NodeMetadata{
			Namespace: "default",
			Labels: map[string]string{
				"istio.io/benchmark": "true",
			},
			IstioVersion: "1.9.0",
		},
		IstioVersion:    &model.IstioVersion{Major: 1, Minor: 8},
		ConfigNamespace: "default",
	}

	t := &testing.T{}
	s := xds.NewFakeDiscoveryServer(t, xds.FakeOptions{
		Configs: configs,
	})

	initPushContext(s.Env(), proxy)

	return s, proxy
}

func initPushContext(env *model.Environment, proxy *model.Proxy) {
	env.PushContext.InitContext(env, nil, nil)
	proxy.SetSidecarScope(env.PushContext)
	proxy.SetGatewaysForProxy(env.PushContext)
	proxy.SetServiceInstances(env.ServiceDiscovery)
}
