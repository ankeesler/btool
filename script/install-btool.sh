#!/bin/bash

set -eo pipefail

curl \
  -L https://github.com/ankeesler/btool/releases/download/0.0/btool-0.0-linux-x86_64.gz \
  -o - \
  | gunzip \
  > /usr/local/bin/btool \
  && chmod +x /usr/local/bin/btool
