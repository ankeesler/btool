#!/bin/bash

set -eo pipefail

if [[ -z "$BTOOL" ]]; then
  BTOOL="$(which btool)"
fi 

cd "$(dirname "${BASH_SOURCE[0]}")/.."
$BTOOL -loglevel error -run -target core/flags_test
$BTOOL -loglevel error -run -target core/cmd_test
$BTOOL -loglevel error -run -target core/err_test

$BTOOL -loglevel error -run -target util/fs/fs_test

$BTOOL -loglevel error -run -target node/node_test
$BTOOL -loglevel error -run -target node/store_test

$BTOOL -loglevel error -run -target app/app_test
$BTOOL -loglevel error -run -target app/builder/builder_test
$BTOOL -loglevel error -run -target app/builder/currenter_impl_test
$BTOOL -loglevel error -run -target app/cleaner/cleaner_test
$BTOOL -loglevel error -run -target app/collector/collector_test
$BTOOL -loglevel error -run -target app/collector/fs/fs_collectini_test
$BTOOL -loglevel error -run -target app/lister/lister_test
$BTOOL -loglevel error -run -target app/runner/runner_test

$BTOOL -loglevel error -run -target btool
