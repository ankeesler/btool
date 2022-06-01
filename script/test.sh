#!/bin/bash

set -eo pipefail

if [[ -z "$BTOOL" ]]; then
  BTOOL="$(which btool)"
fi 

if [[ -z "$REGISTRY" ]]; then
  REGISTRY="https://ankeesler.github.io/btool"
fi 

usage() {
  echo "usage: test.sh [-i] [-u]"
  echo "  -i    Run integration tests    (default: true)"
  echo "  -u    Run unit tests           (default: true)"
}

run_unit_test() {
  echo "running test $1"
  "$BTOOL" -root source -registry "$REGISTRY" -loglevel error -run -target "$1"
}

run_unit_tests() {
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
      run_unit_test "$t"
    fi
  done
}

run_integration_tests() {
  python integration/integration.py \
         "$BTOOL" /tmp/btool-new "$REGISTRY"
  python integration/integration.py \
         /tmp/btool-new /tmp/btool-new-new "$REGISTRY"

  if [[ "$(uname)" == "Linux" ]]; then
    if ! which valgrind; then
      apt-get update && apt-get install valgrind -y
    fi
    valgrind /tmp/btool-new -root source -target btool
  fi
}

cd "$(dirname "${BASH_SOURCE[0]}")/.."

unit_tests=1
integration_tests=1
while getopts hiu o; do
  case "$o" in
    u) unit_tests=1
       integration_tests=0
       ;;
    i) unit_tests=0
       integration_tests=1
       ;;
    h) usage
       exit 0
       ;;

    [?]) usage
         exit 1
         ;;
  esac
done

if [[ "$unit_tests" -eq 1 ]]; then
  run_unit_tests
fi

if [[ "$integration_tests" -eq 1 ]]; then
  run_integration_tests
fi
