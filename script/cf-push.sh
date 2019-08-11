#!/bin/bash

set -eou pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."
tmpdir="$(mktemp -d)"
cf push btool_registry \
   --docker-image ankeesler/btoolregistry \
   --docker-username ankeesler
curl btoolregistry.cfapps.io
rm -rf "$tmpdir"
