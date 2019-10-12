#!/bin/bash

set -eou pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."
docker run --rm -it -v "$PWD:/etc/btool-mount" -w /etc/btool-mount ankeesler/btool ./script/really-test.sh
