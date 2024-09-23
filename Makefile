CURRENTPATH = $(shell echo $(PWD))
WORKDIR = /src/github.com/getyourguide.com/istio-config-validator


run:
	go run cmd/istio-config-validator/main.go -t examples/ examples/

build:
	go build -o istio-config-validator cmd/istio-config-validator/main.go

install:
	go install cmd/istio-config-validator/main.go

test:
	go test -race -count=1 ./...
