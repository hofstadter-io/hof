@experiment(aliasv2)
package veg

import (
	"strings"

	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

bins: multi: [string]: _

dist: {
	[!~"(images|sbom)"]~(k,_): {@env()
		#hof: {id: "dist-\(k)", metadata: {name: string | *id}}
		name: string | *#hof.metadata.name
	}

	meta: env.#ExportDir & {
		path: "dist/meta"
		sources: [
			// root.src.changelog,
			// bins.checksum,

			for sbom in dist.sbom { sbom },
		]
		wipe: true
	}

	cuemod: env.#ExportDir & {
		name: "cue-module"
		path: "dist/cuemod"
		sources: [root.src.cuemod]
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
			#hof: id: string | *"dist-veg-\(k)"
			#hof: metadata: name: #hof.id
			reg: root.flags.registry
			name: "veg-\(k)"
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
			(_f): env.#ExportImage & {
				@env()
				#hof: id: "dist-\(_f)"
				#hof: metadata: name: #hof.id
				image: F.img,
			}
		}
	}

	sbom: {
		cuemod: env.#CuefigSBOM & {
			@env()
			#hof: id: "sbom-cuemod"
			#hof: metadata: name: #hof.id
			path: "cuemod.cue"
			format: "cue"
			data: dist.cuemod
		}
		bins: env.#CuefigSBOM & {
			@env()
			#hof: id: "sbom-bins"
			#hof: metadata: name: #hof.id
			path: "bins.cue"
			format: "cue"
			data: dist.bins
		}
		for i, img in dist.images {
			(i): env.#CuefigSBOM & {
				@env()
				_id: strings.TrimPrefix(img.#hof.id, "dist-")
				#hof: id: "sbom-\(_id)"
				#hof: metadata: name: #hof.id
				path: "\(_id).cue"
				format: "cue"
				data: img
			}
		}
	}

}
