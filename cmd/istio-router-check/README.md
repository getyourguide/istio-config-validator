# Istio Router Check

An _experimental_ wrapper command for Envoy [Route Table Check Tool](https://www.envoyproxy.io/docs/envoy/latest/operations/tools/route_table_check_tool#install-tools-route-table-check-tool).

1. It parses Istio configuration into the Envoy [HTTP Route](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/route/v3/route_components.proto#http-route-components-proto) format required by the tool.
2. It expects `router_check_tool` binary in the PATH. The tool must be built using Istio's filters and still patches needed to work with Istio. It is currently maintained in a fork of istio proxy in [getyourguide/proxy](https://github.com/getyourguide/proxy).


## Usage

```text
Usage:
  istio-router-check [flags]

Flags:
  -v, -- int                        log verbosity level (default 1)
  -c, --config-dir string           directory containing virtualservices
      --covall                      measure coverage by checking all route fields
      --detailed-coverage           print detailed coverage information
  -d, --details                     print detailed information about the test results (default true)
      --disable-deprecation-check   disable deprecation check (default true)
  -f, --fail-under float            threshold for failure
  -h, --help                        help for istio-router-check
      --only-show-failures          only show failures
  -o, --output-dir string           output directory for coverage information
  -t, --test-dir string             directory containing tests
```

### -t \<string>, –test-path \<string>

Path to a tool config JSON file. The tool config JSON file schema is found in config. The tool config input file specifies urls (composed of authorities and paths) and expected route parameter values. Additional parameters such as additional headers are optional.

Schema: All internal schemas in the tool are based on proto3.

### -c \<string>, –config-path \<string>

Path to a router config file (YAML or JSON). The router config file schema is found in config and the config file extension must reflect its file type (for instance, .json for JSON and .yaml for YAML).

### -o \<string>, –output-path \<string>

Path to a file where to write test results as binary proto. If the file already exists, an attempt to overwrite it will be made. The validation result schema is found in proto3.

### -d, –details

Show detailed test execution results. The first line indicates the test name.

### --only-show-failures

Displays test results for failed tests. Omits test names for passing tests if the details flag is set.

### -f, --fail-under

Represents a percent value for route test coverage under which the run should fail.

### --covall

Enables comprehensive code coverage percent calculation taking into account all the possible asserts. Displays missing tests.

### --disable-deprecation-check

Disables the deprecation check for RouteConfiguration proto.

### --detailed-coverage

Enables displaying of not covered routes for non-comprehensive code coverage mode.

### -h, –help

Displays usage information and exits.

## Running

```bash
$ docker run -v $(pwd)/examples:/examples --rm docker.io/getyourguide/istio-router-check:release-1.22 -c /examples/virtualservice.yml -t examples/test.yml

test details.prod.svc.cluster.local/api/v2/products
test details.prod.svc.cluster.local/api/v2/items
Current route coverage: 50%
```
