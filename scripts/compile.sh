#!/usr/bin/env bash

set -x

go test -v && CGO_ENABLED=0 go build -v -a