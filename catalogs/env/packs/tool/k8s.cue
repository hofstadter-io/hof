package tool

import (
	"encoding/yaml"
	"strings"

	"github.com/hofstadter-io/hof/catalogs/env/packs/lang"
	"github.com/hofstadter-io/hof/schemas/env"
)

k8s: kubectl: env.Sh & {
	#ver:  string | *"1.31.0"
	#arch: *"arm64" | "amd64"

	_src: "https://storage.googleapis.com/kubernetes-release/release/v\(#ver)/bin/linux/\(#arch)/kubectl"

	script: """
    set -eou pipefail

    cd /tmp
    wget -q \(_src)
    chmod +x ./kubectl
    mv kubectl /usr/local/bin/kubectl
    rm -rf /tmp/*
    """
}

k8s: helm: env.Sh & {
	#ver:  string | *"4.0.4"
	#arch: *"arm64" | "amd64"

	_src:  "https://get.helm.sh/\(_file)"
	_file: "helm-v\(#ver)-linux-\(#arch).tar.gz"

	script: """
    set -eou pipefail

    cd /tmp
    wget -q \(_src)
    tar -xzf \(_file)
    mv linux-\(#arch)/helm /usr/local/bin/helm
    rm -rf /tmp/*
    """
}

// basically the same as github, but without the version in teh filname
k8s: crane: env.Sh & {
	#ver:    string | *"0.20.7"
	#arch:   string | *"arm64" | "x86_64"
	#distro: string | *"Linux"

	#repo: "google/go-containerregistry"
	#bins: ["crane", "krane"]

	#name: string | *strings.Split(#repo, "/")[1]
	_bins: strings.Join(#bins, " ")

	_file:  "\(#name)_\(#distro)_\(#arch).tar.gz"
	_src:   "https://github.com/\(#repo)/releases/download/v\(#ver)/\(_file)"
	script: """
    set -eou pipefail

    cd /tmp
    wget -q \(_src)
    tar -xzf \(_file)
    mv \(_bins) /usr/local/bin/
    rm -rf /tmp/*
    """
}

k8s: kind: {
	binary: env.File & {
		path: "/usr/local/bin/kind"
		content: lang.go.install.moduleBinary & {
			#params: {
				module:  "sigs.k8s.io/kind"
				version: "v0.31.0"
			}
		}
	}
	config: env.File & {
		path:    "/veg/config/kind-config.yaml"
		content: yaml.Marshal(_config)
	}
	_config: {
		kind:       "Cluster"
		apiVersion: "kind.x-k8s.io/v1alpha4"
		name:       "kind-veg"
		networking: {
			// WARNING: It is _strongly_ recommended that you keep this the default
			// (127.0.0.1) for security reasons. However it is possible to change this.
			apiServerAddress: "127.0.0.1"

			// By default the API server listens on a random open port.
			// You may choose a specific port but probably don't need to in most cases.
			// Using a random port makes it easier to spin up multiple clusters.
			apiServerPort: 6443
		}

		nodes: [{
			role: "control-plane"

			kubeadmConfigPatches: [
				_certSanPatch,
			]
		},

			// extraPortMappings: [{
			// // 	containerPort: 80
			// // 	hostPort: 80
			// // 	protocol: "TCP"
			// // },{
			// 	containerPort: 6443
			// 	hostPort: 6443
			// 	protocol: "TCP"
			// }]
		]

		// nodes: [{
		// 	role: "control-plane"
		// },{
		// 	role: "worker"
		// },{
		// 	role: "worker"
		// }]
	}
}
_certSanPatch: """
	kind: ClusterConfiguration
	apiServer:
	  certSANs:
	    - host.docker.internal
	    - localhost
	    - "127.0.0.1"
	"""
