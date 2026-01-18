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

			for sbom in dist.sbom {sbom},
		]
		wipe: true
	}

	cuemod: env.#ExportDir & {
		name: "cue-module"
		path: "dist/cuemod"
		sources: [root.src.cuemod]
		wipe: true
	}

	github: env.#ExportDir & {
		name: "github-files"
		path: "dist/github"
		sources: [
			for key, val in hof.cli.matrix if key != "name" {val},

		]
		trimPrefix: "./bins"
		wipe:       true
	}

	vscode: env.#ExportDir & {
		name: "vscode-files"
		path: "dist/vscode"
		sources: [
			extn.vscode.vsix,
		]
		wipe: true
	}

	images: {
		[string]~(k,_): env.#PublishImage & {
			@env()
			#hof: id: string | *"dist-veg-\(k)"
			#hof: metadata: name: #hof.id
			reg:  root.flags.registry
			name: string | *"veg-\(k)"
		}
		min: {image: root.ctr.min}
		dev: {image: root.ctr.dev}
		// the dev image with adk/dagger added
		hof: {
			image: env.#Container & {
				from: root.ctr.dev
				steps: [
					env.Dir & {path: "/adk", source: src.adk.fork},
					env.Dir & {path: "/dagger", source: src.dagger.fork},
				]
			}
		}
		ops: {image: root.ctr["ops-all"], name: "veg-ops"}
		for f, F in root.fmtr {
			let _f = "fmt-\(f)"
			(_f): {
				@env()
				#hof: id: "dist-\(_f)"
				#hof: metadata: name: #hof.id
				name: "\(_f)"
				image: F.img
			}
		}
	}

	sbom: {
		for k in ["cuemod", "github", "vscode"] {
			(k): env.#CuefigSBOM & {
				@env()
				#hof: id: "sbom-\(k)"
				#hof: metadata: name: #hof.id
				path:   "\(k).cue"
				format: "cue"
				data:   dist[k]
			}
		}
		for i, img in dist.images {
			(i): env.#CuefigSBOM & {
				@env()
				_id: strings.TrimPrefix(img.#hof.id, "dist-")
				#hof: id: "sbom-\(_id)"
				#hof: metadata: name: #hof.id
				path:   "\(_id).cue"
				format: "cue"
				data:   img
			}
		}
	}

}
