package lang

import (
	"strings"

	"github.com/hofstadter-io/hof/catalogs/env/bases"
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

go: {
	#ver: string | *"1.25.5"

	#goos:   string | *flags.goos
	#goarch: string | *flags.arch

	envSets: {
		default: [
			env.EnvVars & {
				PATH:    "$PATH:/usr/local/go/bin"
				GOBIN:   "/usr/local/bin" // go install to /usr/local/bin
				GOCACHE: "/cache/go"      // intermediate build artifacts
				GOPATH:  "/go"            // mod / pkg / sumdb cache
			},
		]
	}

	caches: {
		goBuild: env.#Cache & {
			name: "go-build-\(#ver)-\(#goos)-\(#goarch)"
		}
		goMods: env.#Cache & {
			name: "go-mods"
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

	defaultSteps: [// new way
		envSets.default,
		mounts.goBuild,
		mounts.goMods,
		install.cli,
		install.devExtras,
		install.lsp,
	]

	install: {
		cli: [
			env.Sh & {
				_file:  "go\(#ver).\(#goos)-\(#goarch).tar.gz"
				_src:   "https://go.dev/dl/\(_file)"
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
			env.Sh & {script: "go install golang.org/x/tools/gopls@latest"},
		]

		moduleBinary: env.#File & {
			#params: {
				module:       string
				version:      string | *"latest"
				_installName: "\(module)@\(version)"
				_parts:       strings.Split(module, "/")
				_name:        _parts[len(_parts)-1]
			}
			path: string | *"/usr/local/bin/\(#params._name)"
			source: env.#Container & {
				from: ctr.base
				steps: [
					env.Sh & {script: "go install \(#params._installName)"},
				]
			}
		}
	}

	ctr: {
		base: env.#Container & {
			@env(pack-lang-go-ctr-base)
			from: bases.debian13.minimal
			steps: [
				utils.apt.install & {#pkgs: [
					"g++",
					"gcc",
					"libc6-dev",
					"netbase",
					"pkg-config",
					"sq",
				]},
				defaultSteps,
			]
		}
		lsp: env.#Container & {
			@env(pack-lang-go-ctr-lsp)
			from: base
			steps: [
				env.Expose & {port: 4000},
			]
		}

		// A dev container with gopls running and attached
		dev: env.#Container & {
			@env(pack-lang-go-ctr-dev)
			from: lsp
			steps: [
				env.BindService & {service: svc.lsp},
			]
		}
	}

	// gopls-as-a-service
	svc: {
		lsp: env.#Service & {
			@env(pack-lang-go-svc-lsp)

			// There is also a built in MCP server!
			#port: int | *0
			name:  "gopls"
			ports: [{name: "lsp", port: 4000, frontend: #port}]
			args: ["gopls", "serve", "-port=4000"]
			// source: _
			source: ctr.lsp & {name: "gopls"}
		}
	}
}
