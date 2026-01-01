@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

bins: multi: [string]: _

dist: {
	[!~"images"]~(k,_): {@env()
		#hof: metadata: {id: "dist-\(k)", name: string | *id}
		name: string | *#hof.metadata.name
	}

	meta: env.#ExportDir & {
		path: "dist/meta"
		sources: [
			root.src.changelog,
			bins.checksum,
			dist.sboms,
		]
		wipe: true
	}

	checksum: env.#File
	sboms: env.#Dir

	cuemod: env.#ExportDir & {
		path: "dist/cuemod"
		sources: [root.src.cuemod]
		wipe: true
	}

	bins: env.#ExportDir & {
		path: "dist/bins"
		sources: [
			root.bins.hof,
			// for k,val in root.bins.multi if k != "name" { val },
		]
		bundlePath: "./bins"
		wipe:       true
	}

	// vscode: env.#ExportDir & {
	//   path: "dist/vscode"
	//   wipe: true
	//   sources: []
	// }

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
