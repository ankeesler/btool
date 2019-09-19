#!/bin/bash

set -eou pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."
docker run --rm -it -v "$PWD:/etc/btool" -w /etc/btool golang:buster go test ./...

go build -o /tmp/btool ./cmd/btool
BTOOL=/tmp/btool ./cc/script/test.sh
