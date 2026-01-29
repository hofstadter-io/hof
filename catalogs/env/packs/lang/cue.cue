package lang

import (
	"github.com/hofstadter-io/hof/catalogs/env/bases"
	// "github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

cue: {
	#ver: string | *"0.15.4"

	#goos:   string | *flags.goos
	#goarch: string | *flags.arch

	caches: {
		cueMods: env.Volume & {
			name: "cue-mods-\(#ver)"
			type: "cache"
		}
	}

	defaultSteps: [
		// any env vars?
		// with cache
		install.cli,
	]

	install: {
		cli: [
			env.Sh & {

				_file:  "cue_v\(#ver)_\(#goos)_\(#goarch).tar.gz"
				_src:   "https://github.com/cue-lang/cue/releases/download/v\(#ver)/\(_file)"
				script: """
					cd /tmp
					wget -q \(_src)
					tar -xzf \(_file)
					mv cue /usr/local/bin/
					rm -rf /tmp/*
					"""
			},
		]
	}

	ctr: {

		// base container with cue binary
		base: env.#Container & {
			@env(pack-lang-cue-ctr-base)
			from: bases.debian13.minimal
			steps: [defaultSteps]
		}

		// container for cue lsp to run in
		lsp: env.#Container & {
			@env(pack-lang-cue-ctr-lsp)
			from: base
			steps: [env.Expose & {port: 4000}]
		}

		// A dev container with cuepls running and attached
		dev: env.#Container & {
			@env(pack-lang-cue-ctr-dev)
			from: lsp
			steps: [
				env.BindService & {service: svc.lsp},
			]
		}
	}

	// cuepls-as-a-service
	svc: {
		lsp: env.#Service & {
			@env(pack-lang-cue-svc-lsp)
			#port: int | *0
			name:  "cuepls"
			ports: [{name: "lsp", port: 4000, frontend: #port}]
			args: ["cue", "lsp", "-port=4000"]
			source: ctr.lsp & {name: "cuepls"}
		}
	}
}
