package veg

import (
	// "github.com/hofstadter-io/hof/catalogs/env/bases"
	"github.com/hofstadter-io/hof/catalogs/env/packs"
	"github.com/hofstadter-io/hof/schemas/env"
)

hack: {
	// used as a simple reproducer to determine that veg-dagger-engine
	// was missing a volume mount and using crazy amounts of disk
	diskUsage: env.#Container & {
		@env(hack-diskUsage)
		from: "ghcr.io/hofstadter-io/veg-hof:v0.7.0-alpha.1"
		steps: [
			env.Mount & {path: "/work", source: src.code},
		]
	}

	sbom: {
		cuefig: env.#CuefigSBOM & {
			@env(hack-cuefig-sbom)
			path:   "hack-cuefig-sbom.cue"
			format: "cue"
			data:   diskUsage
		}
	}

	install: {
		crane: env.#ExportDir & {
			@env(hack-install-crane)
			path: "/usr/local/bin"
			sources: [(packs.containers.crane & {#distro: "Darwin"}).files]
			// sources: [files]
			include: ["crane", "krane"]
			wipe: true
		}
	}
}
