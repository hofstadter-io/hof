package utils

import (
	"strings"

	"github.com/hofstadter-io/hof/schemas/env"
)

dnf: {
	// generalized dnf package install
	install: env.Bash & {
		#pkgs: [...string]
		script: "dnf install -y \(strings.Join(#pkgs, " "))"
	}

	// runs apt-get update, do this once early
	update: env.Bash & {script: "dnf makecache"}

	// You should NEVER need this again!
	// we use caches to do even better than either method
	// 1. same size savings as ( [update -> install -> clean] )
	// 2. save time with cache ( update -> [install] ... magic)
	// anyway, it cleans apt stuff
	clean: env.Bash & {script: "dnf clean all"}
}
