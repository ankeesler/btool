#!/bin/bash

set -eo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."
BTOOL="btool -clean" ./script/test.sh
