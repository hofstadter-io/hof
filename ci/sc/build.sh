#!/bin/bash
set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
ROOT="$DIR/../.."

docker build -t hof/sonar-scanner-cli $ROOT/ci/sc/docker/

