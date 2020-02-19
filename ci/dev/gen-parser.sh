#!/bin/bash
set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
ROOT="$DIR/../.."

PDIR=$ROOT/pkg/lang/parser

pigeon $@ -o $PDIR/hof.go $PDIR/hof.peg

