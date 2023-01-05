#!/bin/bash

env $(cat .env | xargs) air -c .air.toml &
  cd ./dashboard && npm start