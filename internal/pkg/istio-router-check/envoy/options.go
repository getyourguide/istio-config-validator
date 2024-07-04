package envoy

import (
	"istio.io/istio/pkg/config"
)

type optionFunc func(*routeGenerator)

// WithConfigs sets the configs to be used by the route generator.
func WithConfigs(cfgs []config.Config) optionFunc {
	return func(rg *routeGenerator) {
		rg.configs = cfgs
	}
}

// WithGateway generate routes with the provided gateway view. If the gateway does not exist in the provided configs,
// it will be created with the default values accepting "*" as hosts.
func WithGateway(name string) optionFunc {
	return func(rg *routeGenerator) {
		rg.gatewayName = name
	}
}
