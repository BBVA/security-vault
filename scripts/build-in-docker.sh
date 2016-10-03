#!/usr/bin/env bash

set -u # exit when your script tries to use undeclared variables
set -e # make your script exit when a command fails
set -x # trace what gets executed. Useful for debugging

export GO_BUILDER=golang:1.7-wheezy
export REGISTRY_URL="894431722748.dkr.ecr.us-east-1.amazonaws.com"

BASENAME=$(basename $(pwd))
DIRNAME=$(basename $(dirname $(pwd)))

: ${GO_PIPELINE_NAME:="$BASENAME"}
: ${GO_PIPELINE_COUNTER:=dev}
: ${VERSION:="0.1.0"}

BUILD_UID=${BUILD_UID:-`id -u`}
BUILD_GID=${BUILD_GID:-`getent group docker | cut -d: -f3`}

BUILD_COMMAND="docker run  \
       -u $BUILD_UID:$BUILD_GID \
       -e HOME=/tmp \
       -e GOPATH=/go \
       -e GO_PIPELINE_COUNTER \
       -v $(pwd):/go/src/${DIRNAME}/${BASENAME}/ \
       --rm \
       --name=${GO_PIPELINE_NAME}-${GO_PIPELINE_COUNTER} \
       -w /go/src/${DIRNAME}/${BASENAME}/ \
       ${GO_BUILDER} \
       ./scripts/build.sh"

$BUILD_COMMAND

docker build -t ${REGISTRY_URL}/${BASENAME}:${VERSION}-${GO_PIPELINE_COUNTER} .

