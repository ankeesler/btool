#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

outdir=/tmp/btoolregistry

usage() {
  echo "build-registry.sh [h] [-o outdir]"
  echo "  -h           print this message"
  echo "  -o <outdir>  output directory (default: '$outdir')"
}

while getopts ho: o; do
  case "$o" in
    h) usage
       exit
       ;;
    o) outdir="$OPTARG"
       ;;

    [?]) usage
         exit 1
         ;;
  esac
done

mkdir -p "$outdir"

# Copy gaggles
find registry -name "*.yml" -and -not -name "index.yml" -exec cp {} "$outdir" \;

# Create index
cat <<EOF >"$outdir/index.yml"
---
files:
EOF
pushd "$outdir" >/dev/null
  find . -name "*.yml" \
    | sed -e 's#./##'  \
    | sort \
    | xargs -n1 shasum -a 256 \
    | awk '{printf "- path: %s\n  sha256: %s\n", $2, $1}' >> "$outdir/index.yml"
popd >/dev/null
