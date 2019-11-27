#!/bin/bash

set -eo pipefail

if [[ -z "$BTOOL" ]]; then
  BTOOL="$(which btool)"
fi 

if [[ -z "$REGISTRY" ]]; then
  REGISTRY="https://btoolregistry.cfapps.io"
fi 

cd "$(dirname "${BASH_SOURCE[0]}")/.."
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target util/cmd_test
# "$BTOOL" -registry "$REGISTRY" -loglevel error -run -target util/download_test fails in no-network environments
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target util/flags_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target util/fs/fs_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target util/string/string_test
# "$BTOOL" -registry "$REGISTRY" -loglevel error -run -target util/unzip_test fails on linux!

"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target node/node_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target node/property_store_test

"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/app_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/builder/builder_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/builder/currenter_impl_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/cleaner/cleaner_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/collector/base_collectini_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/collector/collector_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/collector/store_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/collector/cc/exe_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/collector/cc/inc_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/collector/cc/includes_parser_impl_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/collector/cc/obj_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/collector/cc/resolver_factory_delegate_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/collector/fs/fs_collectini_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/collector/registry/registry_collectini_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/collector/registry/gaggle_collector_impl_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/collector/registry/resolver_factory_delegate_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/collector/registry/fs_registry_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/collector/registry/yaml_serializer_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/lister/lister_test
"$BTOOL" -registry "$REGISTRY" -loglevel error -run -target app/runner/runner_test

python integration/integration.py \
       "$BTOOL" /tmp/btool-from-go "$REGISTRY"
python integration/integration.py \
       /tmp/btool-from-go /tmp/btool-from-cc "$REGISTRY"
python integration/integration.py \
       /tmp/btool-from-cc /tmp/btool-from-cc-from-cc "$REGISTRY"
