#!/usr/bin/env bash

set -x # trace what gets executed. Useful for debugging
set -e # make your script exit when a command fails
set -u # exit when your script tries to use undeclared variables


source ./common.sh

# Deploy / update  stack to rancher
docker-rancher-tools rancher-compose -p ${STACK_NAME} -r target/service/rancher-compose.yml -f target/service/docker-compose.yml up -u -c -d