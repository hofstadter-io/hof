package util

import (
	"strings"

	"github.com/hofstadter-io/hof/schemas/env"
)

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
