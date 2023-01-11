#!/bin/bash

env $(cat .env | xargs) go run -tags ee ./cmd/hatchet-server-migrate/main.go