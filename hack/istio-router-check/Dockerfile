FROM 130607246975.dkr.ecr.eu-central-1.amazonaws.com/base/golang-1.23:build-1.0.0 AS build

WORKDIR /work

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin/istio-router-check hack/istio-router-check/main.go


FROM getyourguide/router-check-tool:release-1.22

COPY --from=build /work/bin/istio-router-check /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/istio-router-check"]
