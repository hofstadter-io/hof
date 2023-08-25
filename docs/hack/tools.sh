#!/bin/bash
set -eou pipefail

HUGO_VER="0.111.3"
CUE_VER="v0.6.0"

mkdir tmp
pushd tmp

wget https://github.com/gohugoio/hugo/releases/download/v${HUGO_VER}/hugo_extended_${HUGO_VER}_Linux-64bit.tar.gz -O hugo.tar.gz
tar -xf hugo.tar.gz
chmod +x hugo
sudo mv hugo /usr/local/bin/hugo

wget https://github.com/cue-lang/cue/releases/download/${CUE_VER}/cue_${CUE_VER}_linux_amd64.tar.gz -O cue.tar.gz
tar -xf cue.tar.gz
sudo mv cue /usr/local/bin/cue

popd
rm -rf tmp
