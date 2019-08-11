#!/bin/bash

set -eou pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."
docker build -t ankeesler/btoolregistry .
