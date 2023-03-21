# Base Go environment
# -------------------
FROM golang:1.19-alpine as base
WORKDIR /hatchet

RUN apk update && apk add --no-cache wget gcc musl-dev git protoc

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

# install terraform
RUN wget https://releases.hashicorp.com/terraform/1.1.8/terraform_1.1.8_linux_amd64.zip
RUN unzip terraform_1.1.8_linux_amd64.zip && rm terraform_1.1.8_linux_amd64.zip

# Go build environment
# --------------------
FROM base AS build-go

ARG version=v0.0.1

# build proto files
RUN sh ./hack/proto/proto.sh

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=$GOPATH/pkg/mod \
    go build -ldflags="-w -s -X 'main.Version=${version}'" -a -o ./bin/hatchet-runner ./cmd/hatchet-runner

# Deployment environment
# ----------------------
FROM alpine
RUN apk update

COPY --from=base /hatchet/terraform /usr/bin/terraform
COPY --from=build-go /hatchet/bin/hatchet-runner /hatchet/

CMD /hatchet/hatchet-runner
