#!/bin/bash
set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
ROOT="$DIR/../.."

pigeon $@ \
  -o $ROOT/pkg/parser/hof.go \
  $ROOT/pkg/parser/hof.peg

