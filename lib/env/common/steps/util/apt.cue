package util

import (
	"strings"

	"github.com/hofstadter-io/hof/schemas/env"
)

apt: {
	caches: {
		varLib: env.#Cache & {
			name: "debian-13-var-lib-cache"
		}
	}

	mounts: {
		varLib: env.Mount & {
			path:   "/var/lib/apt/lists"
			source: apt.caches.varLib
		}
	}

	// generalized apt package install
	install: env.Bash & {
		#pkgs: [...string]
		script: "apt-get install -y --no-install-recommends \(strings.Join(#pkgs, " "))"
	}

	// runs apt-get update, do this once early
	update: env.Bash & {script: "apt-get update -y"}

	// You should NEVER need this again!
	// we use caches to do even better than either method
	// 1. same size savings as ( [update -> install -> clean] )
	// 2. save time with cache ( update -> [install] ... magic)
	// anyway, it cleans apt stuff
	clean: env.Bash & {
		script: """
			apt-get dist-clean
			rm -rf /var/lib/apt/lists/*
			"""
	}
}
