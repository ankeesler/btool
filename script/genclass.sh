#!/bin/bash

set -eo pipefail

if [[ -z "$1" ]]; then
  echo "usage: genclass.sh path/to/cc/class"
  exit 1
fi

path="$1"
if [[ "$path" == "-" ]]; then
  path="some/path/to/some_class"
  h="/dev/stdout"
  cc="/dev/stdout"
  testcc="/dev/stdout"
else
  h="source/${path}.h"
  cc="source/${path}.cc"
  testcc="source/${path}_test.cc"
fi

# gen h
ifndef="BTOOL_"
ifndef="${ifndef}$(echo "$path" | tr "[:lower:]" "[:upper:]" | tr / _)"
ifndef="${ifndef}_H_"

namespace="btool::"
namespace="${namespace}$(dirname "$path" | sed -e 's#/#::#g')"

class=
upper="1"
for c in $(basename "$path" | sed -e 's/\(.\)/\1 /g'); do
  if [[ "$c" == "_" ]]; then
    upper="1"
    continue
  fi

  if [[ "$upper" -eq "1" ]]; then
    class="${class}$(echo "$c" | tr "[:lower:]" "[:upper:]")"
  else
    class="${class}$c"
  fi

  upper="0"
done

cat <<EOF >"$h"
#ifndef $ifndef
#define $ifndef

namespace $namespace {

class $class {
};

};  // namespace $namespace

#endif  // $ifndef
EOF
echo "wrote $h"

# gen cc
cat <<EOF >"$cc"
#include "${path}.h"

namespace $namespace {

};  // namespace $namespace
EOF
echo "wrote $cc"

# gen testcc
cat <<EOF >"$testcc"
#include "${path}.h"

#include "gtest/gtest.h"

TEST($class, Yeah) {
}
EOF
echo "wrote $testcc"
