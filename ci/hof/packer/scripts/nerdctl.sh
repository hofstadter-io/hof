#!/bin/bash
set -euo pipefail

echo "nerdctl"

VER="1.4.0"
ARCH=${ARCH:-amd}

# nerdctl full
pushd /tmp
wget -q https://github.com/containerd/nerdctl/releases/download/v${VER}/nerdctl-full-${VER}-linux-${ARCH}64.tar.gz
sudo tar Cxzf /usr/local nerdctl-full-${VER}-linux-${ARCH}64.tar.gz
popd

# enable containerd for rootfull
sudo systemctl enable --now containerd

# deal with apparmor
sudo nerdctl apparmor load

# not a huge fan of this, but get error without that
# only workaround I have found thus far
# https://github.com/containerd/nerdctl/discussions/1536
sudo ln -s /usr/sbin/iptables /usr/local/bin/iptables
sudo ln -s /usr/sbin/ip6tables /usr/local/bin/ip6tables

# validate
sudo nerdctl version
sudo nerdctl info
sudo nerdctl run hello-world

# sudoless has to be setup per user by linking files, so we do that later
