#!/bin/bash

set -eo pipefail

if [[ -z "$BTOOL" ]]; then
  BTOOL="$(which btool)"
fi 

if [[ -z "$REGISTRY" ]]; then
  REGISTRY="https://btoolregistry.cfapps.io"
fi 

run_test() {
  valgrind "$BTOOL" -root source -registry "$REGISTRY" -loglevel error -run -target "$1"
}

cd "$(dirname "${BASH_SOURCE[0]}")/.."

tests="$(find . -name "*_test.cc" | sed -e 's!./source/!!;s!\.cc!!')"
for t in $tests; do
  if [[ "$t" == "util/unzip_test" ]]; then
    echo
    echo "**********"
    echo
    echo "skipping util/unzip_test - libzip does not work on linux"
    echo
    echo "**********"
    echo
  else
    run_test "$t"
  fi
done

python integration/integration.py \
       "$BTOOL" /tmp/btool-new "$REGISTRY"
python integration/integration.py \
       /tmp/btool-new /tmp/btool-new-new "$REGISTRY"

if which valgrind; then
  valgrind /tmp/btool-new -root source -target btool
fi
