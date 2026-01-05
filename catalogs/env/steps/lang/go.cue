package lang

import (
	"strings"

	"github.com/hofstadter-io/hof/catalogs/env/bases"
	"github.com/hofstadter-io/hof/schemas/env"
)

// versions
go: {
	#ver: string | *"1.25.5"

	#arch: *"arm64" | "amd64"
	// os, always assume linux

	envSets: {
		default: [
			env.EnvVar & {PATH: "$PATH:/usr/local/go/bin"},
			env.EnvVar & {GOBIN: "/usr/local/bin"}, // go install to /usr/local/bin
			env.EnvVar & {GOCACHE: "/cache/go"},    // intermediate build artifacts
			env.EnvVar & {GOPATH: "/go"},           // mod / pkg / sumdb cache
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

	setupSteps: [
		envSets.default,
		mounts.goBuild,
		mounts.goMods,
		install.cli,
	]

	default: defaultSteps
	defaultSteps: [// new way
		envSets.default,
		mounts.goBuild,
		mounts.goMods,
		install.cli,
		install.devExtras,
		install.lsp,
	]

	dev: env.#Container & {
		from: bases.debian13.default
		steps: defaultSteps
	}

	install: {
		cli: [
			env.Sh & {
				_file:   "go\(#ver).linux-\(#arch).tar.gz"
				_src:    "https://go.dev/dl/\(_file)"
				script: """
					cd /tmp
					wget -q \(_src)
					tar -C /usr/local -xzf \(_file)
					rm -rf /tmp/*
					"""
			},
		]

		devExtras: [
			env.Sh & {
				script: """
					# lint tools
					go install honnef.co/go/tools/cmd/staticcheck@latest
					go install github.com/mgechev/revive@latest
					curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b /usr/local/bin v2.7.2
					"""
			},
		]

		lsp: [
			env.Sh & { script: "go install golang.org/x/tools/gopls@latest" },
		]

		moduleBinary: env.#File & {
			#params: {
				module: string
				version: string | *"latest"
				_installName: "\(module)@\(version)"
				_parts: strings.Split(module,"/")
				_name: _parts[len(_parts)-1]
			}
			path: string | *"/usr/local/bin/\(#params._name)"
			source: env.#Container & {
				from: go.dev
				steps: [
					env.Sh & { script: "go install \(#params._installName)"}
				]
			}
		}
	}
}
