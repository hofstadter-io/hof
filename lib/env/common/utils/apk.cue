package utils

import (
	"strings"

	"github.com/hofstadter-io/hof/schemas/env"
)

apk: {
	// generalized dnf package install
	install: env.Sh & {
		#pkgs: [...string]
		script: "apk add \(strings.Join(#pkgs, " "))"
	}

	// runs apt-get update, do this once early
	update: env.Sh & {script: "apk update --no-interactive"}
	upgrade: env.Sh & {script: "apk upgrade --no-interactive"}

	// You should NEVER need this again!
	// we use caches to do even better than either method
	// 1. same size savings as ( [update -> install -> clean] )
	// 2. save time with cache ( update -> [install] ... magic)
	// anyway, it cleans apt stuff
	// clean: env.Bash & {script: "dnf clean all"}
}
