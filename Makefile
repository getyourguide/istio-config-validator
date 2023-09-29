CURRENTPATH = $(shell echo $(PWD))
WORKDIR = /src/github.com/getyourguide.com/istio-config-validator
GOLINT_CMD=golangci-lint
GO_FILES=$(shell find . -type f -regex ".*go")

.PHONY: run lint fix

run:
	docker run -it --rm --name istio_config_validator \
				-v ${CURRENTPATH}:${WORKDIR} \
				-w ${WORKDIR} \
				golang:1.21 \
				go run cmd/istio-config-validator/main.go -t examples/ examples/

define check_linter
	@if ! command -v $(GOLINT_CMD) > /dev/null; then\
		echo "$(GOLINT_CMD) must be installed to run linting: try brew install $(GOLINT_CMD)";\
		exit 1;\
	fi
endef

lint: .golangci.yml ## Lint the code

fix: ## Fix any issues found during linting
	@$(call check_linter)
	$(GOLINT_CMD) run --fix

.golangci.yml: $(GO_FILES)
	@$(call check_linter)
	$(GOLINT_CMD) version
	$(GOLINT_CMD) run
	touch .golangci.yml