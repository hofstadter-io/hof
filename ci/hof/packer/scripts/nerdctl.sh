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
#sudo systemctl enable --now containerd

# validate
#sudo nerdctl version
#sudo nerdctl info
#sudo nerdctl run hello-world

#
# OR
#

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
