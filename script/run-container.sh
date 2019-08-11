#!/bin/bash

set -eou pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."
docker run --rm -it -p 8080:8080 ankeesler/btoolregistry "$@"
