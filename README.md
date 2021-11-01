# istio-config-validator

[![Go Report Card](https://goreportcard.com/badge/github.com/getyourguide/istio-config-validator)](https://goreportcard.com/report/github.com/getyourguide/istio-config-validator)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/6bee3a704e8648949523cdcfcefacc1f)](https://www.codacy.com?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=getyourguide/istio-config-validator&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/6bee3a704e8648949523cdcfcefacc1f)](https://www.codacy.com?utm_source=github.com&utm_medium=referral&utm_content=getyourguide/istio-config-validator&utm_campaign=Badge_Coverage)

> The `istio-config-validator` tool is a **Work In Progress** project.

It provides to developers and cluster operators a way to test their changes in VirtualServices. We do it by mocking Istio/Envoy behavior to decide to which destination the request would go to. Eg:

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: example
  namespace: example
spec:
  hosts:
    - www.example.com
    - example.com
  http:
    - match:
        - uri:
            regex: /users(/.*)?
          headers:
            x-user-type:
              exact: qa
      route:
        - destination:
            host: users.users.svc.cluster.local
            port:
              number: 80
    - route:
        - destination:
            host: monolith.monolith.svc.cluster.local
```

Given the above `VirtualService`, developers can introduce test cases that covers the intended behavior as the following:

```yaml
testCases:
  - description: Only QA users should go to the new Users microservice.
    wantMatch: true
    request:
      authority: ["www.example.com", "example.com"]
      method: ["GET", "OPTIONS", "POST"]
      uri: ["/users", "/users/"]
      headers:
        x-user-type: qa
    route:
    - destination:
        host: users.users.svc.cluster.local
        port:
          number: 80
  - description: Fallback other user types to the monolith
    wantMatch: true
    request:
      authority: ["www.example.com", "example.com"]
      method: ["GET", "OPTIONS", "POST"]
      uri: ["/users", "/users/"]
    route:
    - destination:
        host: monolith.monolith.svc.cluster.local
```

Have a look in the [TestCase Reference](docs/test-cases.md) to learn more how to define the tests.

## Installation
Either install the go package
```
# go install github.com/getyourguide/istio-config-validator/cmd/istio-config-validator@latest
```
Or alternatively install the docker image
```
# docker pull getyourguide/istio-config-validator:latest
```

## Usage

```
# istio-config-validator -h
Usage: istio-config-validator -t <testcases1.yml|testcasesdir1> [-t <testcases2.yml|testcasesdir2> ...] <istioconfig1.yml|istioconfigdir1> [<istioconfig2.yml|istioconfigdir2> ...]

  -t value
        Testcase files/folders
```

```
# istio-config-validator -t examples/virtualservice_test.yml examples/virtualservice.yml
2020-05-29T18:45:39.261018Z     info    running test: happy path users
2020-05-29T18:45:39.261106Z     info    PASS input:[{www.example.com GET /users map[x-user-id:abc123]}]
2020-05-29T18:45:39.261128Z     info    PASS input:[{www.example.com GET /users/ map[x-user-id:abc123]}]
2020-05-29T18:45:39.261141Z     info    PASS input:[{www.example.com POST /users map[x-user-id:abc123]}]
2020-05-29T18:45:39.261157Z     info    PASS input:[{www.example.com POST /users/ map[x-user-id:abc123]}]
2020-05-29T18:45:39.261169Z     info    PASS input:[{example.com GET /users map[x-user-id:abc123]}]
2020-05-29T18:45:39.261184Z     info    PASS input:[{example.com GET /users/ map[x-user-id:abc123]}]
2020-05-29T18:45:39.261207Z     info    PASS input:[{example.com POST /users map[x-user-id:abc123]}]
2020-05-29T18:45:39.261220Z     info    PASS input:[{example.com POST /users/ map[x-user-id:abc123]}]
===========================
2020-05-29T18:45:39.261228Z     info    running test: Partner service only accepts GET or OPTIONS
2020-05-29T18:45:39.261256Z     info    PASS input:[{example.com PUT /partners map[]}]
2020-05-29T18:45:39.261274Z     info    PASS input:[{example.com PUT /partners/1 map[]}]
2020-05-29T18:45:39.261284Z     info    PASS input:[{example.com POST /partners map[]}]
2020-05-29T18:45:39.261900Z     info    PASS input:[{example.com POST /partners/1 map[]}]
2020-05-29T18:45:39.261940Z     info    PASS input:[{example.com PATCH /partners map[]}]
2020-05-29T18:45:39.261984Z     info    PASS input:[{example.com PATCH /partners/1 map[]}]
===========================
```

## Contributing

If you're interested in contributing to this project or running a dev version, have a look into the [CONTRIBUTING](CONTRIBUTING.md) document

## Known Limitations

The API for test cases does not cover all aspects of VirtualServices.

-   Supported [HTTPMatchRequests](https://istio.io/docs/reference/config/networking/virtual-service/#HTTPMatchRequest) fields to match requests against are: `authority`, `method`, `headers` and `uri`.
    -   Not supported ones: `scheme`, `port`, `queryParams`, etc.

-   Supported assert against [HTTPRouteDestination](https://istio.io/docs/reference/config/networking/virtual-service/#HTTPRouteDestination) and [HTTPRewrite](https://istio.io/docs/reference/config/networking/virtual-service/#HTTPRewrite)
    -   Not supported ones: [HTTPRedirect](https://istio.io/docs/reference/config/networking/virtual-service/#HTTPRedirect), etc.

## Security

For sensitive security matters please contact [security@getyourguide.com](mailto:security@getyourguide.com).


## Legal

Copyright 2020 GetYourGuide GmbH.

istio-config-validator is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full text.
