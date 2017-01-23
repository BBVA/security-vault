#!/usr/bin/env bash

export CONTAINER_NAME="security-vault"                  # Name of the container
export SERVICE_NAME="${CONTAINER_NAME}"      # Name of the service, advised to match project folder
export SOURCE_PACKAGE="descinet.bbva.es"                # Package where the project resides in the GOPATH
export AT_DOCKER_NAME="security-vault-at"               # Acceptance Tests container name defined in ../deploy/acceptance-test/docker-compose.yml
export VAULT_CONFIGURATOR_DOCKER_NAME="vault-configurator"
export VERSION="0.2.0"

export REGISTRY_URL="#PLACEHOLDER#"                     # Docker registry url
export STACK_NAME="security-vault"                      # Rancher stack name
export RANCHER_CATALOG_URI="#PLACEHOLDER#"              # Rancher catalog url
export GO_BUILDER=golang:1.7-wheezy                     # Container used to build the application

export GO_PIPELINE_COUNTER=${GO_PIPELINE_COUNTER:-dev}  # Build number from GOCD pipeline, can be substituted by any build number from any CI tool

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