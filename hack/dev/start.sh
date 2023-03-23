#!/bin/bash

set -a
. .env
set +a

air -c .air.toml & cd ./dashboard && npm start