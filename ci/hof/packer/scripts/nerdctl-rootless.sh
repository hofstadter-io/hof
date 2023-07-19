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

# enable rootless mode (also remove sudo below)
sudo apt-get install -y \
	dbus-user-session \
	uidmap
systemctl --user start dbus
containerd-rootless-setuptool.sh install

# not a huge fan of this, but get error without that
# only workaround I have found thus far
# https://github.com/containerd/nerdctl/discussions/1536
sudo ln -s /usr/sbin/iptables /usr/local/bin/iptables
sudo ln -s /usr/sbin/ip6tables /usr/local/bin/ip6tables

# validate
nerdctl version
nerdctl info
nerdctl run hello-world
