#!/bin/bash

env $(cat .env | xargs) go run -tags ee ./cmd/hatchet-temporal