@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

bins: multi: [string]: _

dist: {
	[!~"images"]~(k,_): {@env()
		#hof: { id: "dist-\(k)", metadata: {name: string | *id}}
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
		path: "dist/cuemod"
		sources: [src.code]
		include: [
			"cue.mod/module.cue",
			// "*.cue", // eventually, when we rework all of ci, use .veg more, and have a root index that imports many things, like a mega package if the user wants
			"schemas",
			"examples",
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
		path: "dist/bins"
		sources: [
			root.bins.hof,
			for key,val in root.bins.matrix if key != "name" { val },
		]
		// maybe this is better as trimPrefix or extractPath, this name is not clear
		bundlePath: "./bins"
		wipe:       true
	}

	images: {
		[string]~(k,_): {
			@env()
			name: "veg-\(k)"
		}
		min: env.#ExportImage & {image: root.ctr.min}
		dev: env.#ExportImage & {image: root.ctr.dev}
		for f, F in root.fmtr {
			"fmtr-\(f)": env.#ExportImage & {image: F.img}
		}
	}

}
