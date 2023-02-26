#!/bin/bash
set -eou pipefail

# setup the development workspace

HUGO_VER="0.109.0"

npm install
ln -s ../node_modules assets/node_modules
npm install broken-link-checker -g
mkdir tmp

pushd tmp

wget https://github.com/gohugoio/hugo/releases/download/v${HUGO_VER}/hugo_extended_${HUGO_VER}_Linux-64bit.tar.gz -O hugo.tar.gz
tar -xf hugo.tar.gz
chmod +x hugo
sudo mv hugo /usr/local/bin/hugo

popd
rm -rf tmp
