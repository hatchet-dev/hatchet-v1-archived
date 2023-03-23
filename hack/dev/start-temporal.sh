#!/bin/bash

set -a
. .env
set +a

go run -tags ee ./cmd/hatchet-temporal