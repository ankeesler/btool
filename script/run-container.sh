#!/bin/bash

set -eo pipefail

image="$1"
if [[ -z "$image" ]]; then
    echo "usage: run-container.sh <image>"
    exit 1
fi
shift

cd "$(dirname "${BASH_SOURCE[0]}")/.."
docker run -v "$PWD:/etc/btool-mount" -w /etc/btool-mount --rm -it -p 80:80 "ankeesler/$image" "$@"
