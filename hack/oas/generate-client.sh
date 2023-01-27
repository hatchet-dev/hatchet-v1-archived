#!/bin/bash

swagger-codegen generate -i ./bin/oas/openapi.yaml -l go -o ./api/v1/client/swagger