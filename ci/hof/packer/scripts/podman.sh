#!/bin/bash
set -euo pipefail

echo "podman"

sudo apt-get install -y podman

podman version
podman info

podman run hello-world
