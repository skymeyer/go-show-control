#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

go generate ./...
go build -a -o bin/show-control main.go
