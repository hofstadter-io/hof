package utils

import (
	"strings"

	"github.com/hofstadter-io/hof/schemas/env"
)

apk: {
	install: env.Sh & {
		#pkgs: [string, ...string]
		script: "apk add --no-cache \(strings.Join(#pkgs, " "))"
	}
	upgrade: env.Sh & {
		#pkgs: [...string]
		_pkgs: string | *""
		if len(#pkgs) > 0 {
			_pkgs: strings.Join(#pkgs, " ")
		}
		script: "apk upgrade --no-cache --no-interactive \(_pkgs)"
	}
	update: env.Sh & {script: "apk update --no-cache --no-interactive"}
	clean: env.Sh & {script: "apk cache clean && apk cache purge"}
}

apt: {
	install: env.Sh & {
		#pkgs: [string, ...string]
		script: "apt-get install -y --no-install-recommends \(strings.Join(#pkgs, " "))"
	}
	upgrade: env.Sh & {
		#pkgs: [...string]
		_pkgs: string | *""
		if len(#pkgs) > 0 {
			_pkgs: strings.Join(#pkgs, " ")
		}
		script: "apt-get upgrade -y --no-install-recommends \(_pkgs)"
	}
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

dnf: {
	install: env.Sh & {
		#pkgs: [string, ...string]
		script: "dnf install -y --nodocs \(strings.Join(#pkgs, " "))"
	}
	upgrade: env.Sh & {
		#pkgs: [...string]
		_pkgs: string | *""
		if len(#pkgs) > 0 {
			_pkgs: strings.Join(#pkgs, " ")
		}
		script: "dnf upgrade -y --nodocs \(_pkgs)"
	}
	update: env.Sh & {script: "dnf makecache"}
	clean: env.Sh & {script: "dnf clean all"}
}
