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
for f in $(find "$outdir" -name "*.yml" | sort); do
  cat <<EOF >>"$outdir/index.yml"
- path: $(basename "$f")
  sha256: $(cat "$f" | openssl sha256 -hex | cut -f 2 -d ' ')
EOF
done
