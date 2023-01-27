BINDIR      := $(CURDIR)/bin
VERSION ?= dev

test:
	go test ./... -tags test -cover

migrate:
	sh ./hack/dev/migrate.sh

start-dev:
	sh ./hack/dev/start.sh

build-openapi:
	sh ./hack/oas/generate-spec.sh

build-proto:
	sh ./hack/proto/proto.sh

build-api-client: build-openapi
	cd dashboard; npx swagger-typescript-api -p ../bin/oas/openapi.yaml -o ./src/shared/api/generated -n hatchet.ts --modular

build-api-client-golang: build-openapi
	sh ./hack/oas/generate-client.sh
