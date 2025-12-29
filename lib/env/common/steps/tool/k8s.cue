package tool

import (
	"strings"

	"github.com/hofstadter-io/hof/schemas/env"
)

k8s: kubectl: env.Exec & {
	#ver:  string | *"1.31.0"
	#arch: *"arm64" | "amd64"

	args: ["sh", "-c", _script]

	_src: "https://storage.googleapis.com/kubernetes-release/release/v\(#ver)/bin/linux/\(#arch)/kubectl"

	_script: """
    set -eou pipefail

    cd /tmp
    wget -q \(_src)
    chmod +x ./kubectl
    mv kubectl /usr/local/bin/kubectl
    rm -rf /tmp/*
    """
}

k8s: helm: env.Exec & {
	#ver:  string | *"4.0.4"
	#arch: *"arm64" | "amd64"

	args: ["sh", "-c", _script]

	_src:  "https://get.helm.sh/\(_file)"
	_file: "helm-v\(#ver)-linux-\(#arch).tar.gz"

	_script: """
    set -eou pipefail

    cd /tmp
    wget -q \(_src)
    tar -xzf \(_file)
    mv linux-\(#arch)/helm /usr/local/bin/helm
    rm -rf /tmp/*
    """
}

// basically the same as github, but without the version in teh filname
k8s: crane: env.Exec & {
	#ver:    string | *"0.20.7"
	#arch:   string | *"arm64" | "x86_64"
	#distro: string | *"Linux"

	#repo: "google/go-containerregistry"
	#bins: ["crane", "krane"]

	#name: string | *strings.Split(#repo, "/")[1]
	_bins: strings.Join(#bins, " ")

	args: ["sh", "-c", _script]

	_file:   "\(#name)_\(#distro)_\(#arch).tar.gz"
	_src:    "https://github.com/\(#repo)/releases/download/v\(#ver)/\(_file)"
	_script: """
    set -eou pipefail

    cd /tmp
    wget -q \(_src)
    tar -xzf \(_file)
    mv \(_bins) /usr/local/bin/
    rm -rf /tmp/*
    """
}
