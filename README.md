# istio-config-validator - test your istio config
[![Go Report Card](https://goreportcard.com/badge/github.com/getyourguide.com/istio-config-validator)](https://goreportcard.com/report/github.com/getyourguide.com/istio-config-validator)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/6bee3a704e8648949523cdcfcefacc1f)](https://www.codacy.com?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=getyourguide/istio-config-validator&amp;utm_campaign=Badge_Grade)

## Installation
Either install the go package
```
# go get -u github.com/getyourguide/istio-config-validator/cmd/istio-config-validator
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

## Development

Compilation and building is handled in the Docker container:
-   checkout the git repo
-   in the repo folder, run `docker build -t istio-config-validator:latest .`
