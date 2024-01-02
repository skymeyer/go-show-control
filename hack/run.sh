#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

go run main.go run --debug
