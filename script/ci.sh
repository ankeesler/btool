#!/bin/bash

set -eou pipefail

announce() {
  echo "========================================================================"
  echo "          $1"  
  echo "========================================================================"
  "$@"
}

cd "$(dirname "${BASH_SOURCE[0]}")/.."
announce ./script/test.sh
announce ./script/build-container.sh btool
announce ./script/build-container.sh btoolregistry
announce ./script/push-container.sh btool
announce ./script/push-container.sh btoolregistry
announce ./script/cf-push.sh
