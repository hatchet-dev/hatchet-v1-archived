#!/bin/bash
#
# Builds auto-generated protobuf files

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/v1/server/pb/provisioner.proto