#!/bin/bash

set -eo pipefail

image="$1"
if [[ -z "$image" ]]; then
    echo "usage: run-container.sh <image>"
    exit 1
fi
shift

cd "$(dirname "${BASH_SOURCE[0]}")/.."
docker run --rm -it -p 8080:8080 "ankeesler/$image" "$@"
