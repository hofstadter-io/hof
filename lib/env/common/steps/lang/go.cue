package lang

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

// versions
go: {
	#ver: string | *"1.25.4"

	#arch: *"arm64" | "amd64"
	// os, always assume linux

	envSets: {
		default: [
			env.Env & {PATH: "$PATH:/usr/local/go/bin"},
			env.Env & {GOBIN: "/usr/local/bin"}, // go install to /usr/local/bin
			env.Env & {GOCACHE: "/cache/go"},    // intermediate build artifacts
			env.Env & {GOPATH: "/go"},           // mod / pkg / sumdb cache
		]
	}

	caches: {
		goBuild: env.#Cache & {
			name: "go-build-\(#ver)-\(#arch)"
		}
		goMods: env.#Cache & {
			name: "go-mods-\(#ver)-\(#arch)"
		}
	}

	mounts: {
		goBuild: env.Mount & {
			path:   "/cache/go"
			source: caches.goBuild
		}
		goMods: env.Mount & {
			path:   "/go"
			source: caches.goMods
		}
	}

	default: [
		envSets.default,
		mounts.goBuild,
		mounts.goMods,
		install.cli,
		install.devExtras,
		install.lsp,
	]

	install: {

		cli: [
			env.Exec & {
				args: ["sh", "-c", _script]

				_file:   "go\(#ver).linux-\(#arch).tar.gz"
				_src:    "https://go.dev/dl/\(_file)"
				_script: """
					set -eou pipefail

					cd /tmp
					wget -q \(_src)
					tar -C /usr/local -xzf \(_file)
					rm -rf /tmp/*
					"""
			},
		]

		devExtras: [
			env.Exec & {
				args: ["sh", "-c", _script]

				_script: """
					# lint tools
					go install honnef.co/go/tools/cmd/staticcheck@latest
					go install github.com/mgechev/revive@latest
					curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b /usr/local/bin v2.7.2
					"""
			},
		]

		lsp: [
			env.Exec & {
				args: ["sh", "-c", _script]

				_script: """
					# LSP
					go install golang.org/x/tools/gopls@latest
					"""
			},
		]
	}
}
