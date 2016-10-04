#!/usr/bin/env bash

export SERVICE_NAME="cloudframe-security-vault"
export SOURCE_PACKAGE="descinet.bbva.es"
export AT_DOCKER_NAME="cloudframe-security-vault-at"  # defined in ../deploy/acceptance-test/docker-compose.yml
export VERSION="0.1.0"

export REGISTRY_URL="894431722748.dkr.ecr.us-east-1.amazonaws.com"
export STACK_NAME="cloudframe-security-vault"
export RANCHER_CATALOG_URI="https://descinet.bbva.es/stash/scm/cloudframe/rancher-catalog-security.git"
export GO_BUILDER=golang:1.7-wheezy

export GO_PIPELINE_COUNTER=${GO_PIPELINE_COUNTER:-dev}

function getVersion {
    if [ ! -f "../version" ]; then
        echo "${VERSION}-${GO_PIPELINE_COUNTER}" > ../version;
    fi
    cat ../version
}

function docker-rancher-tools {
    docker run -i --rm \
    -v "$PWD/../:/app" -w /app \
    -e RANCHER_URL \
    -e RANCHER_SECRET_KEY \
    -e RANCHER_ACCESS_KEY \
    -e JAVA_OPTS \
    -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY \
    -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_KEY \
    894431722748.dkr.ecr.us-east-1.amazonaws.com/rancher-tools:0.0.1-17 $@
}

function docker-rancher-api-cli {
    docker run -i --rm \
    -e RANCHER_URL -e RANCHER_ACCESS_KEY -e RANCHER_SECRET_KEY \
    894431722748.dkr.ecr.us-east-1.amazonaws.com/rancher-api-cli:0.0.1-4 "$@"
}

function is_aws_instance {
    curl -s --max-time 1 http://169.254.169.254/latest/meta-data/
}

if is_aws_instance; then aws=true; else aws=false; fi