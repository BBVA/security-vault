#!/usr/bin/env bash

set -u # exit when your script tries to use undeclared variables
set -e # make your script exit when a command fails
set -x # trace what gets executed. Useful for debugging

source ./common.sh

VERSION=$(getVersion)

BUILD_UID=${BUILD_UID:-`id -u`}
BUILD_GID=${BUILD_GID:-`getent group docker | cut -d: -f3`}

docker run  \
       -u ${BUILD_UID}:${BUILD_GID} \
       -e HOME=/tmp \
       -e GOPATH=/go \
       -v $(pwd)/..:/go/src/${SOURCE_PACKAGE}/${SERVICE_NAME}/ \
       --rm \
       -w /go/src/${SOURCE_PACKAGE}/${SERVICE_NAME}/ \
       ${GO_BUILDER} \
       ./scripts/compile.sh

docker build --force-rm -f  ../Dockerfile -t ${REGISTRY_URL}/${SERVICE_NAME}:${VERSION} ..
docker build --force-rm -f ../Dockerfile_at -t ${REGISTRY_URL}/${AT_DOCKER_NAME}:${VERSION} ..

if $aws ; then
    docker push ${REGISTRY_URL}/${SERVICE_NAME}:${VERSION}
    docker push ${REGISTRY_URL}/${AT_DOCKER_NAME}:${VERSION}
fi
