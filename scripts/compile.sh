#!/usr/bin/env bash

set -x
set -u
set -e

pushd test
    go test -v
popd

CGO_ENABLED=0 go build -v -a