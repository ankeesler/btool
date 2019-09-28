#!/bin/bash

set -eo pipefail

if [[ -z "$BTOOL" ]]; then
  BTOOL="$(which btool)"
fi 

cd "$(dirname "${BASH_SOURCE[0]}")/.."
"$BTOOL" -loglevel error -run -target node/node_test
"$BTOOL" -loglevel error -run -target core/flags_test
"$BTOOL" -loglevel error -run -target app/lister/lister_test
"$BTOOL" -loglevel error -run -target btool
