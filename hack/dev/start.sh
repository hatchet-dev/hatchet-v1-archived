#!/bin/bash

env $(cat .env | xargs) go run ./cmd/hatchet-server/main.go