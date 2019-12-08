#!/bin/bash

set -euo pipefail

usage() {
  echo "gennodes.sh <-u url> <-h header_dir> <-s source_dir> <-l library> [-i include_dir] [-d]"
  echo
  echo "example: gennodes.sh \\"
  echo "  -u https://github.com/OpenSSL_1_1_1d.tar.gz \\"
  echo "  -h include \\"
  echo "  -s crypto \\"
  echo "  -l libcrypto.a"
}

debug="0"
debug() {
  if [[ "$debug" == "1" ]]; then
    echo debug: $@ 1>&2
  fi
}

url=
header_dir=
source_dir=
library=
include_dir=
while getopts "u:h:s:l:i:d" o; do
  case "$o" in
    u) url="$OPTARG"
       ;;
    h) header_dir="$OPTARG"
       ;;
    s) source_dir="$OPTARG"
       ;;
    l) library="$OPTARG"
       ;;
    i) include_dir="$OPTARG"
       ;;
    d) debug="1"
       ;;
    [?]) usage
         exit 1
         ;;
  esac
done

if [[ -z "$url" ]] || [[ -z "$header_dir" ]] || [[ -z "$source_dir" ]] || [[ -z "$library" ]]; then
  usage
  exit 1
fi

dir="$(mktemp -d)"
archive_name="$(basename $url)"
archive="$dir/$archive_name"
curl -L -o "$archive" "$url"
sha256="$(shasum -a 256 <"$archive" | awk '{print $1}')"
tar xzf "$archive" -C "$dir"
debug "archive: $archive"

root="$dir/$(ls -1 "$dir" | grep -v "$archive_name")"
debug "root: $root"

headers="$(find $root/$header_dir -name "*.h" | sed -e "s!$dir/!!g")"
debug "found $(echo "$headers" | wc -w) headers"

sources="$(find $root/$source_dir -name "*.c" | sed -e "s!$dir/!!g")"
debug "found $(echo "$sources" | wc -w) sources"

cat <<EOF
nodes:
# Archives
- name: $archive_name
  dependencies:
    - \$this
  resolver:
    name: download
    config:
      url: $url
      sha256: $sha256

EOF

echo "# Headers"
for header in $headers; do
  cat <<EOF
- name: $header
  dependencies:
    - $archive_name
  labels:
    io.btool.collector.cc.libraries:
      - $library
  resolver:
    name: untar

EOF
done

echo "# Sources"
for source in $sources; do
  cat <<EOF
- name: $source
  dependencies:
    - $archive_name
  labels:
    io.btool.collector.cc.includePaths:
      - $(basename $root)/$header_dir
EOF

  if [[ ! -z "$include_dir" ]]; then
    echo "      - $(basename $root)/$include_dir"
  fi

  cat << EOF
  resolver:
    name: untar

EOF
done

objects=
echo "# Objects"
for source in $sources; do
  object="$(echo "$source" | sed -e 's/\.c$/.o/g')"
  objects="$objects $object"

  cat <<EOF
- name: $object
  dependencies:
    - $source
  resolver:
    name: compileC

EOF
done

echo "# Libraries"
echo "- name: $library"
echo "  dependencies:"
for object in $objects; do
  echo "    - $object"
done
echo "  resolver:"
echo "    name: archive"

rm -rf "$dir"
