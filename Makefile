BINDIR      := $(CURDIR)/bin
VERSION ?= dev

test:
	go test ./... -tags test -cover

start-dev:
	sh ./hack/dev/start.sh

build-openapi:
	sh ./hack/oas/generate-spec.sh

build-api-client: build-openapi
	cd dashboard; npx swagger-typescript-api -p ../bin/oas/openapi.yaml -o ./src/shared/api -n hatchet.ts --modular