# https://suva.sh/posts/well-documented-makefiles/#simple-makefile
.DEFAULT_GOAL:=help
SHELL:=/bin/bash

.PHONY: help test build push plan apply destroy

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

test: ## Run tests
	cd test && go test

build: ## Build docker image for docker-compose
	cd hello-world-docker && docker build . -t go-webapp

test-unit:  ## Run unit tests
	@echo "Running unit tests..."
	$(eval testunitargs += "-timeout=60m")
	@mkdir -p tmp
	@if which gotestsum > /dev/null 2>&1 ; then \
		echo "gotestsum --format=standard-verbose --junitfile=tmp/unit-report.xml --" $(testunitargs); \
		cd test && gotestsum --format=standard-verbose --junitfile=../tmp/unit-report.xml -- $(testunitargs); \
	else \
		echo "go test -v" $(testunitargs); \
		cd test && go test -v $(testunitargs); \
	fi;
