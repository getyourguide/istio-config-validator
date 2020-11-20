module github.com/getyourguide/istio-config-validator

go 1.15

replace github.com/spf13/viper => github.com/istio/viper v1.3.3-0.20190515210538-2789fed3109c

// Old version had no license
replace github.com/chzyer/logex => github.com/chzyer/logex v1.1.11-0.20170329064859-445be9e134b2

// Avoid pulling in incompatible libraries
replace github.com/docker/distribution => github.com/docker/distribution v2.7.1+incompatible

// Avoid pulling in kubernetes/kubernetes
replace github.com/Microsoft/hcsshim => github.com/Microsoft/hcsshim v0.8.8-0.20200421182805-c3e488f0d815

// Client-go does not handle different versions of mergo due to some breaking changes - use the matching version
replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.5

// See https://github.com/kubernetes/kubernetes/issues/92867, there is a bug in the library
replace github.com/evanphx/json-patch => github.com/evanphx/json-patch v0.0.0-20190815234213-e83c0a1c26c8

require (
	github.com/ghodss/yaml v1.0.0
	github.com/golang/protobuf v1.4.3
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.16.0
	istio.io/api v0.0.0-20201120175956-c2df7c41fd8e
	istio.io/client-go v1.8.0
	istio.io/istio v0.0.0-20201120181132-0d30e764cc51
	istio.io/pkg v0.0.0-20201119191759-d9a2706ea471
	k8s.io/apimachinery v0.19.4
)
