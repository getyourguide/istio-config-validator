# Istio Router Check

An _experimental_ wrapper command for Envoy [Route Table Check Tool](https://www.envoyproxy.io/docs/envoy/latest/configuration/operations/tools/router_check).

1. It parses Istio configuration and generate it in the Envoy [HTTP Route](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/route/v3/route_components.proto#http-route-components-proto) format required by the tool.
2. It expects `router_check_tool` binary in the PATH. The tool must be built using Istio's filters and still patches needed to work with Istio. It is currently maintained in a fork of istio proxy in [getyourguide/proxy](https://github.com/getyourguide/proxy).

## Running

```bash
$ docker run -v $(pwd)/examples:/examples --rm docker.io/getyourguide/istio-router-check:release-1.22 -c /examples/virtualservice.yml -t examples/test.yml
test details.prod.svc.cluster.local/api/v2/products
test details.prod.svc.cluster.local/api/v2/items
Current route coverage: 50%
```
