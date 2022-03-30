CURRENTPATH = $(shell echo $(PWD))
WORKDIR = /src/github.com/getyourguide.com/istio-config-validator


run:
	docker run -it --rm --name istio_config_validator \
				-v ${CURRENTPATH}:${WORKDIR} \
				-w ${WORKDIR} \
				golang:1.18 \
				go run cmd/istio-config-validator/main.go -t examples/ examples/
