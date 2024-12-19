# Istio Router Check

An _experimental_ command that generates configs and tests to Envoy [Route Table Check Tool](https://www.envoyproxy.io/docs/envoy/latest/operations/tools/route_table_check_tool#install-tools-route-table-check-tool).

1. It parses Istio configuration such as VirtualServices and outputs Envoy [HTTP Route](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/route/v3/route_components.proto#http-route-components-proto) format required by router check tool.
2. It parses mutitple Envoy Tests and consolidate them into a single file to be used by router check tool.
3. It parses istio-config-validator test format and converts to the router check tool format. Note that this is highly experimental and does not cover all tests.

## Running

Generate routes and consolidate tests:

```shell
$ docker run \
  -v $(pwd)/examples:/examples \
  --rm docker.io/getyourguide/istio-router-check:release-1.24 --config-dir /examples/virtualservices --test-dir examples/envoy-tests/ -o examples/build/

time=2024-07-09T13:41:43.137Z level=INFO msg="reading tests" dir=examples/envoy-tests/
time=2024-07-09T13:41:43.145Z level=INFO msg="writing tests" file=examples/build/tests.json
time=2024-07-09T13:41:43.146Z level=INFO msg="reading virtualservices"
time=2024-07-09T13:41:43.244Z level=INFO msg="writing route" route=80 file=examples/build/route_sidecar_80.json
```

Run envoy route table check tool using the generated files.

```shell
docker run \
  -v $(pwd)/examples:/examples \
  --entrypoint=/usr/local/bin/router_check_tool \
  --rm docker.io/getyourguide/istio-router-check:release-1.24 \
  -c /examples/build/route_sidecar_80.json -t /examples/build/tests.json --only-show-failures --disable-deprecation-check

Current route coverage: 50%
```

## Usage

```shell
docker run \
  --rm docker.io/getyourguide/istio-router-check:release-1.24 -h

Usage:
  istio-router-check [flags]

Flags:
  -v, -- int                Log verbosity level
  -c, --config-dir string   Directory with Istio VirtualService and Gateway files
      --gateway string      Only consider VirtualServices bound to this gateway (i.e: istio-system/istio-ingressgateway)
  -h, --help                Help for istio-router-check
  -o, --output-dir string   Directory to output Envoy routes and tests
  -t, --test-dir string     Directory with Envoy test files
```

Envoy Router Check

```shell
docker run \
  --entrypoint=/usr/local/bin/router_check_tool
  --rm docker.io/getyourguide/istio-router-check:release-1.24 -h


USAGE:

   /usr/local/bin/router_check_tool  [--detailed-coverage] [-o <string>]
                                     [-t <string>] [-c <string>] [--covall]
                                     [-f <float>]
                                     [--disable-deprecation-check]
                                     [--only-show-failures] [-d] [--]
                                     [--version] [-h]
                                     <unlabelledConfigStrings> ...


Where:

   --detailed-coverage
     Show detailed coverage with routes without tests

   -o <string>,  --output-path <string>
     Path to output file to write test results

   -t <string>,  --test-path <string>
     Path to test file.

   -c <string>,  --config-path <string>
     Path to configuration file.

   --covall
     Measure coverage by checking all route fields

   -f <float>,  --fail-under <float>
     Fail if test coverage is under a specified amount

   --disable-deprecation-check
     Disable deprecated fields check

   --only-show-failures
     Only display failing tests

   -d,  --details
     Show detailed test execution results

   --,  --ignore_rest
     Ignores the rest of the labeled arguments following this flag.

   --version
     Displays version information and exits.

   -h,  --help
     Displays usage information and exits.

   <unlabelledConfigStrings>  (accepted multiple times)
     unlabelled configs


   router_check_tool
```
