#!/bin/bash
set -euo pipefail

echo "packages"

sudo apt-get update
sudo apt-get install -y \
	ca-certificates \
	curl \
	git \
	gnupg \
	tree \
	wget

