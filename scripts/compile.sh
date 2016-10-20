#!/usr/bin/env bash

set -x

pushd test
    go test -v
popd

CGO_ENABLED=0 go build -v -a