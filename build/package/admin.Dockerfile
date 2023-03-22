# Base Go environment
# -------------------
FROM golang:1.19-alpine as base
WORKDIR /hatchet

RUN apk update && apk add --no-cache gcc musl-dev git protoc

COPY go.mod go.sum ./
COPY /api ./api
COPY /cmd ./cmd
COPY /ee ./ee
COPY /internal ./internal
COPY /hack ./hack

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

RUN --mount=type=cache,target=$GOPATH/pkg/mod \
    go mod download

# Go build environment
# --------------------
FROM base AS build-go

ARG version=v0.0.1

# build proto files
RUN sh ./hack/proto/proto.sh

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=$GOPATH/pkg/mod \
    go build -ldflags="-w -s -X 'github.com/hatchet-dev/hatchet/cmd/hatchet-admin/cli.Version=${version}'" -a -o ./bin/hatchet-admin ./cmd/hatchet-admin 

# Deployment environment
# ----------------------
FROM alpine
RUN apk update && apk add --no-cache openssl

COPY --from=build-go /hatchet/bin/hatchet-admin /hatchet/

CMD /hatchet/hatchet-admin
