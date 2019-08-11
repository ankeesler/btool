#!/bin/bash

set -eou pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."
tmpdir="$(mktemp -d)"
GOOS=linux go build -o "$tmpdir/registry" github.com/ankeesler/btool/cmd/registry
cp -r data "$tmpdir/data"
cf push btool_registry \
   -c './registry -loglevel trace -dir data -address :$PORT' \
   -b binary_buildpack \
   -p "$tmpdir"
curl btoolregistry.cfapps.io
rm -rf "$tmpdir"
