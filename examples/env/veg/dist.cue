@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

bins: multi: [string]: _

dist: {
	[!~"images"]~(k,_): {@env()
		#hof: {id: "dist-\(k)", metadata: {name: string | *id}}
		name: string | *#hof.metadata.name
	}

	// meta: env.#ExportDir & {
	// 	path: "dist/meta"
	// 	sources: [
	// 		root.src.changelog,
	// 		bins.checksum,
	// 		dist.sboms,
	// 	]
	// 	wipe: true
	// }

	cuemod: env.#ExportDir & {
		name: "cue-module"
		path: "dist/cuemod"
		sources: [src.code]
		include: [
			"cue.mod/module.cue",
			// "*.cue", // eventually, when we rework all of ci, use .veg more, and have a root index that imports many things, like a mega package if the user wants
			"schemas",
			"catalogs/env",
			"examples/env",
			"flow/tasks/*.cue",
			"flow/tasks/*/*.cue",
			"lib/env/common",
			"SECURITY.md",
			"README.md",
			"AGENTS.md",
			"LICENSE",
		]
		wipe: true
	}

	bins: env.#ExportDir & {
		name: "gh-release"
		path: "dist/bins"
		sources: [
			hof.cli.local,
			for key, val in hof.cli.matrix if key != "name" {val}
		]
		// maybe this is better as trimPrefix or extractPath, this name is not clear
		bundlePath: "./bins"
		wipe:       true
	}

	images: {
		[string]~(k,_): {
			@env()
			name: string | *"veg-\(k)"
			reg: "ghcr.io/hofstadter-io"
		}
		min: env.#ExportImage & {image: root.ctr.min}
		dev: env.#ExportImage & {image: root.ctr.dev}
		// the dev image with adk/dagger added
		hof: env.#ExportImage & {
			image: env.#Container & {
				from: root.ctr.dev
				steps: [
					env.Dir & { path: "/adk", source: src.adk.fork },
					env.Dir & { path: "/dagger", source: src.dagger.fork },
				]
			}
		}
		ops: env.#ExportImage & {image: root.ctr["ops-all"], name: "veg-ops"}
		for f, F in root.fmtr {
			let _f = "fmt-\(f)"
			(_f): env.#ExportImage & {image: F.img, name: _f}
		}
	}

}
