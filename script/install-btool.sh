#!/bin/bash

set -eo pipefail

usage() {
  echo "usage: install-btool.sh local        (1)"
  echo "usage: install-btool.sh latest       (2)"
  echo "usage: install-btool.sh <version>    (3)"
  echo
  echo "(1) build locally (requires local btool binary"
  echo "(2) download latest released btool binary"
  echo "(3) download released btool binary with verison <version>"
}

download_and_install() {
  version="$1"
  os="$2"
  curl \
    -L "https://github.com/ankeesler/btool/releases/download/$version/btool-$version-$os-x86_64.gz" \
    -o - \
    -f \
    | gunzip \
    > /usr/local/bin/btool \
    && chmod +x /usr/local/bin/btool
}

version="$1"
if [[ -z "$version" ]]; then
  usage
  exit 1
fi

os="linux"
if [[ "$(uname)" == "Darwin" ]]; then
  os="macos"
fi

echo "install-btool.sh: installing version $version"
if [[ "$version" == "local" ]]; then
  if ! which btool; then
    echo "install-btool.sh: error: no local btool binary"
    exit 1
  fi

  btool -root source -target btool -clean # TODO: remove me when builds work across os's!
  btool -root source -target btool
  cp source/btool /usr/local/bin/btool
elif [[ "$version" == "latest" ]]; then
  version="$(curl -f https://api.github.com/repos/ankeesler/btool/releases/latest | jq -r .tag_name)"
  download_and_install "$version" "$os"
else
  download_and_install "$version" "$os"
fi

