FROM 890546402176.dkr.ecr.eu-central-1.amazonaws.com/hubdockercom/golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/istio-config-validator/
COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o /go/bin/istio-config-validator ./cmd/istio-config-validator/

FROM 890546402176.dkr.ecr.eu-central-1.amazonaws.com/hubdockercom/busybox
COPY --from=builder /go/bin/istio-config-validator /go/bin/istio-config-validator

ENTRYPOINT ["/go/bin/istio-config-validator"]
