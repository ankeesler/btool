#!/bin/bash

set -eo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

etags $(find . -name "*.cc" -or -name "*.h")
