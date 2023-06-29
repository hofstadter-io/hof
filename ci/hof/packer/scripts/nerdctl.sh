#!/bin/bash
set -euo pipefail

echo "nerdctl"

VER="1.4.0"

# nerdctl full
pushd /tmp
wget https://github.com/containerd/nerdctl/releases/download/v${VER}/nerdctl-full-${VER}-linux-amd64.tar.gz
sudo tar Cxzf /usr/local nerdctl-full-${VER}-linux-amd64.tar.gz
popd

# enable containerd for rootfull
sudo systemctl enable --now containerd

# validate
sudo nerdctl version
sudo nerdctl info
sudo nerdctl run hello-world
