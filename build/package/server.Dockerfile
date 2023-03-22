ARG VERSION=v0.0.1

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

ARG VERSION

# build proto files
RUN sh ./hack/proto/proto.sh

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=$GOPATH/pkg/mod \
    go build -ldflags="-w -s -X 'main.Version=${VERSION}'" -a -o ./bin/hatchet-server ./cmd/hatchet-server 

# Webpack build environment
# -------------------------
FROM node:16 as build-webpack
WORKDIR /webpack

ARG NPM_TOKEN
ENV NPM_TOKEN=$NPM_TOKEN

COPY ./dashboard ./

RUN npm install -g npm@8.1

RUN npm i --legacy-peer-deps

ENV NODE_ENV=production

RUN npm run build

# Deployment environment
# ----------------------
FROM alpine
RUN apk update

COPY --from=build-go /hatchet/bin/hatchet-server /hatchet/
COPY --from=build-webpack /webpack/build /hatchet/static

EXPOSE 8080
CMD /hatchet/hatchet-server
