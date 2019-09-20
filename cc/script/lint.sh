#!/bin/bash

set -eo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."
files="$(find . -name "*.h" -or -name "*.cc")"
for file in $files; do
  clang-format -style=Google "$file" > /tmp/btool-lint-tmp
  mv /tmp/btool-lint-tmp "$file"
  echo "$file"
done
