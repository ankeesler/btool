#!/bin/bash

set -eo pipefail

VERSION="0.0"
OS="linux"

if [[ "$(uname)" == "Darwin" ]]; then
  OS="macos"
fi

curl \
  -L "https://github.com/ankeesler/btool/releases/download/$VERSION/btool-$VERSION-$OS-x86_64.gz" \
  -o - \
  | gunzip \
  > /usr/local/bin/btool \
  && chmod +x /usr/local/bin/btool
