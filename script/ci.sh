#!/bin/bash

set -eou pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."
./script/test.sh
./script/build-container.sh
./script/push-container.sh
./script/cf-push.sh
