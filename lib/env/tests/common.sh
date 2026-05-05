#!/usr/bin/env bash
set -euo pipefail

GITROOT=$(git rev-parse --show-toplevel)
SELFDIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# TODO, testscript the shit out of this file

# test base images
pushd $SELFDIR/../common/bases
hof def > /dev/null
hof eval -c=false > /dev/null
hof vet  -c=false > /dev/null
hof env list  -P report > reports/common/bases/hof-env-list.txt
hov env build -P report > reports/common/bases/hof-env-build.txt