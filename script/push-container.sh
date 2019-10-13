#!/bin/bash

set -eo pipefail

image="$1"
if [[ -z "$image" ]]; then
    echo "usage: push-container.sh <image>"
    exit 1
fi

cd "$(dirname "${BASH_SOURCE[0]}")/.."
docker push "ankeesler/$image"
