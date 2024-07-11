FROM golang:1.22-bullseye as builder

WORKDIR /work

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin/istio-router-check hack/istio-router-check/main.go

FROM getyourguide/router-check-tool:release-1.22

COPY --from=builder /work/bin/istio-router-check /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/istio-router-check"]
