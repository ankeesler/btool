#!/bin/bash

set -eou pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."
./script/test.sh
./script/build-container.sh btool
./script/build-container.sh btoolregistry
./script/push-container.sh btool
./script/push-container.sh btoolregistry
./script/cf-push.sh
