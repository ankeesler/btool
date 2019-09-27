#!/bin/bash

set -eo pipefail

if [[ -z "$BTOOL" ]]; then
  BTOOL="$(which btool)"
fi 

cd "$(dirname "${BASH_SOURCE[0]}")/.."
"$BTOOL" -run -target node/node_test
"$BTOOL" -run -target core/flags_test
"$BTOOL" -run -target app/lister/lister_test
"$BTOOL" -run -target btool
