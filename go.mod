module github.com/getyourguide/istio-config-validator

go 1.15

require (
	github.com/ghodss/yaml v1.0.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.16.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
	istio.io/api v0.0.0-20200518162646-2c1705ec4d0b
	istio.io/client-go v0.0.0-20200518164621-ef682e2929e5
	istio.io/pkg v0.0.0-20200526141228-3772d4c49765
	k8s.io/apimachinery v0.20.2
)
