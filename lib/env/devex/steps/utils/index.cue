package utils

import (
	"strings"

	"github.com/hofstadter-io/hof/schemas/env"
)

apt: {
	caches: {
		varLib: env.#Cache & {
			name: "/var/lib/apt/lists-debian-13"
		}
		varCache: env.#Cache & {
			name: "/var/cache/apt-debian-13"
		}
	}

	mounts: {
		varLib: env.Mount & {
			path:   "/var/lib/apt/lists"
			source: apt.caches.varLib
		}
		varCache: env.Mount & {
			path:   "/var/cache/apt/archives"
			source: apt.caches.varCache
		}
	}

	// generalized apt package install
	install: env.Exec & {
		#pkgs: [...string]
		_script: """
      set -eux
      apt-get install -y --no-install-recommends \(strings.Join(#pkgs, " "))
      """
		args: ["bash", "-c", _script]
	}
	// runs apt-get update, do this once early
	update: env.Exec & {
		#pkgs: [...string]
		_script: """
			set -eux
			apt-get update -y
			"""
		args: ["bash", "-c", _script]
	}

	// You should NEVER need this again!
	// we use caches to do even better than either method
	// 1. same size savings as ( [update -> install -> clean] )
	// 2. save time with cache ( update -> [install] ... magic)
	// anyway, it cleans apt stuff
	clean: env.Exec & {
		#pkgs: [...string]
		_script: """
			set -eux
			apt-get dist-clean
			rm -rf /var/lib/apt/lists/*
			"""
		args: ["bash", "-c", _script]
	}
}

// helpers for downloading prebuild binaries
githubBin: env.Exec & {
	#ver: string

	#arch:   string | *"arm64" | "amd64" // todo, this should default to current OS
	#distro: string | *"linux"

	#repo: string
	#name: string | *strings.Split(#repo, "/")[1]
	#bins: [...string] | *[#name]
	_bins: strings.Join(#bins, " ")

	args: ["sh", "-c", _script]

	_file:   "\(#name)_v\(#ver)_\(#distro)_\(#arch).tar.gz"
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
