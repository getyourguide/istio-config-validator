package envoy

import (
	"fmt"
	"io/ioutil"
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
	for _, c := range s.Clusters(proxy) {
		fmt.Printf("clusterName:%s\n", c.GetName())
	}

	// To determine which routes to generate, first gen listeners once (not part of benchmark) and extract routes
	l := s.Discovery.ConfigGenerator.BuildListeners(proxy, s.PushContext())
	for _, r := range s.Routes(proxy) {
		fmt.Printf("route:%s\n", r.GetName())
	}
	routeNames := xdstest.ExtractRoutesFromListeners(l)
	if len(routeNames) == 0 {
		log.Fatal("Got no route names!")
	}
	c := s.Discovery.Generators[v3.RouteType].Generate(proxy, s.PushContext(), &model.WatchedResource{ResourceNames: routeNames}, nil)
	for i, r := range c {
		s, err := (&jsonpb.Marshaler{Indent: "  "}).MarshalToString(r)
		if err != nil {
			log.Fatal(err)
		}

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
		// TODO: if you update this, make sure telemetry.yaml is also updated
		IstioVersion:    &model.IstioVersion{Major: 1, Minor: 6},
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
