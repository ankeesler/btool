#!/bin/bash

set -eou pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."
go vet ./...
go test ./...
go build -o /tmp/btool ./cmd/btool
BTOOL=/tmp/btool ./cc/script/test.sh
