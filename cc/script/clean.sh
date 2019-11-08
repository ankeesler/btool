#!/bin/bash

set -eo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."
find . -name "*.o" | xargs rm -rf
find . -name "*_test" | xargs rm -rf
rm -rf "btool"
