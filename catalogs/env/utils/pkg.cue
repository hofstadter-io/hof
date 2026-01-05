package utils

import (
	"strings"

	"github.com/hofstadter-io/hof/schemas/env"
)

dnf: {
	install: env.Sh & { #pkgs: [...string], script: "dnf install -y --nodocs \(strings.Join(#pkgs, " "))" }
	upgrade: env.Sh & { #pkgs: [...string], script: "dnf upgrade -y --nodocs \(strings.Join(#pkgs, " "))" }
	update: env.Sh & {script: "dnf makecache"}
	clean: env.Sh & {script: "dnf clean all"}
}

apk: {
	install: env.Sh & { #pkgs: [...string], script: "apk add --no-cache \(strings.Join(#pkgs, " "))" }
	upgrade: env.Sh & { #pkgs: [...string], script: "apk upgrade --no-cache --no-interactive \(strings.Join(#pkgs, " "))" }
	update: env.Sh & {script: "apk update --no-cache --no-interactive"}
	clean: env.Sh & {script: "apk cache clean && apk cache purge"}
}

apt: {
	install: env.Sh & { #pkgs: [...string], script: "apt-get install -y --no-install-recommends \(strings.Join(#pkgs, " "))" }
	upgrade: env.Sh & { #pkgs: [...string], script: "apt-get upgrade -y --no-install-recommends \(strings.Join(#pkgs, " "))" }
	update: env.Sh & {script: "apt-get update -y"}
	clean: env.Sh & {
		script: """
			apt-get clean
			apt-get autoclean
			apt-get dist-clean
			rm -rf /var/lib/apt/lists/*
			"""
	}
}