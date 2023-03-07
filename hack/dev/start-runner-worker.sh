#!/bin/bash

env $(cat worker.env | xargs) go run ./cmd/hatchet-runner-worker 