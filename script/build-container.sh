#!/bin/bash

set -eou pipefail

image="$1"
if [[ -z "$image" ]]; then
    echo "usage: build-container.sh <image>"
    exit 1
fi

cd "$(dirname "${BASH_SOURCE[0]}")/.."
docker build -t "ankeesler/$image" -f "image/$image/Dockerfile" .
