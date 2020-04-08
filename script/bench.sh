#!/bin/bash

set -eo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

metadata() {
  source="$1"
  cache="$2"
  cat <<EOF
---
metadata:
  source: $source
  cache: $cache
data:
EOF
}

bench() {
  name="$1"
  threads="$2"
  dry_run="$3"
  echo "- name: $name"
  echo "  threads: $threads"

  if [[ "$dry_run" != "true" ]]; then
    btool -cache "$cache" -root source -target btool -threads "$threads" 2>&1 | tail -2
  fi
}

build_fs() {
  root="$1"
}

clean() {
  cache="$1"
  rm -rf ${cache}/*
  ./script/clean.sh
}

usage() {
  echo "bench.sh [-h] [-n]"
  echo "  -n    Dry run (don't run btool)"
  echo "  -h    Print this message"
}

dry_run="false"
while getopts hn o; do
  case "$o" in
    h) usage
       exit 0
       ;;
    n) dry_run="true"
       ;;
    [?]) usage
         exit 1
         ;;
  esac
done

tmp="$(mktemp -d)"
source="$tmp/source"
cache="$tmp/cache"

metadata "$source" "$cache"
build_fs "$source"

clean "$cache"
bench "deep" 1 "$dry_run"
clean "$cache"
bench "deep" 2 "$dry_run"
clean "$cache"
bench "deep" 3 "$dry_run"

rm -rf "$tmp"
