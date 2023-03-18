#!/bin/bash
set -eou pipefail

# setup the development workspace

npm ci
ln -s ../node_modules assets/node_modules
npm install broken-link-checker -g

