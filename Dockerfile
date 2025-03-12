FROM 130607246975.dkr.ecr.eu-central-1.amazonaws.com/base/golang-1.24:build-1.0.0 AS builder
RUN apt-get update && apt-get install -y git
WORKDIR $GOPATH/src/istio-config-validator/
COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o /go/bin/istio-config-validator ./cmd/istio-config-validator/

FROM 130607246975.dkr.ecr.eu-central-1.amazonaws.com/base/golang-1.24:runtime-1.0.0
COPY --from=builder /go/bin/istio-config-validator /go/bin/istio-config-validator

ENTRYPOINT ["/go/bin/istio-config-validator"]
