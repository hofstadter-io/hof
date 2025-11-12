#!/bin/bash
set -euo pipefail

VERSION=${1:?"Usage: $0 <version>"}

GIT_ROOT=$(git rev-parse --show-toplevel)

rm -rf $GIT_ROOT/.tmp/publish-cue-module
mkdir -p $GIT_ROOT/.tmp/publish-cue-module/{,cue.mod}

include=(
  LICENSE
  README.md
  schemas/*
)

for file in "${include[@]}"; do
  cp -rf "$GIT_ROOT/$file" "$GIT_ROOT/.tmp/publish-cue-module/"
done

cat <<EOF > $GIT_ROOT/.tmp/publish-cue-module/cue.mod/module.cue
module: "github.com/hofstadter-io/schemas"
language: {
	version: "v0.13.0"
}
source: {
	kind: "self"
}
EOF

pushd "$GIT_ROOT/.tmp/publish-cue-module"

tree
# cue mod publish $VERSION

popd