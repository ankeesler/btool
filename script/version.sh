#!/bin/bash

set -eo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."
go run ./cmd/btool/main.go -version | awk '{print $NF}'
