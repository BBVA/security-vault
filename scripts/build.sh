#!/usr/bin/env bash

set -x

go test && CGO_ENABLED=0 go build -v -a