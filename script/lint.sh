#!/bin/bash

set -eo pipefail

install_clang_format_brew() {
  brew install clang-format
}

install_clang_format_linux() {
  apt-get update && apt-get install clang-format -y
}

install_clang_format() {
    case "$(uname)" in
      Darwin)
        install_clang_format_darwin
        ;;
      Linux)
        install_clang_format_linux
        ;;
    esac
}

cd "$(dirname "${BASH_SOURCE[0]}")/.."

if ! which clang-format >/dev/null; then
  install_clang_format
fi

files="$(find . -name "*.h" -or -name "*.cc")"
diff_exists=0
for file in $files; do
  clang-format -style=Google "$file" > /tmp/btool-lint-tmp
  if ! diff /tmp/btool-lint-tmp "$file" >/dev/null; then
     diff_exists=1
  fi
  mv /tmp/btool-lint-tmp "$file"
  echo "$file"
done

exit "$diff_exists"
