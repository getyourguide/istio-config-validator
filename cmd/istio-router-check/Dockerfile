FROM golang:1.22-bullseye as builder

WORKDIR /work

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin/istio-router-check cmd/istio-router-check/main.go

# FIXME: Use GetYourGuide registry
FROM cainelli/router_check_tool:release-1.22

COPY --from=builder /work/bin/istio-router-check /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/istio-router-check"]
