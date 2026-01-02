@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/lib/env/common/packs/lang"
	"github.com/hofstadter-io/hof/schemas/env"
)

// go multi-target builds
// from source build too for local setup

ctr: {
	builder: env.#Container & {
		@env(), @id(hof-cli-builder)
		from: lang.go.ctr.base
		steps: [
			env.Dir & {path: "/adk", source: src.adk},
			env.Dir & {path: "/dagger", source: src.dagger},
			env.Dir & {path: "/work", source: src.code},
		]
	}
	built: env.#Container & {
		@env(), @id(hof-cli-built)
		from: builder
		steps: [
			env.EnvVar & {GOOS: flags.goos, GOARCH: flags.arch},
			env.Exec & {args: ["go", "build", "-o", "./bins/hof", "./cmd/hof"]},
		]
	}
}

hof: cli: env.File & {path: "/usr/local/bin/hof", content: bins.matrix["linux-arm64"]}

bins: {
	[string]~(k,_): {name: "bin-\(k)"}
	hof: env.#File & {@env(), path: "./bins/hof", source: ctr.built}
	matrix: {
		_goos: ["linux", "darwin"]
		_arch: ["amd64", "arm64"]
		for _g in _goos for _a in _arch
		// for _g in _goos for _a in _arch {
		// 	let _short = "\(_g)-\(_a)"
		// 	let _bin = "hof-\(_short)"
		// 	"\(_short)": env.#File & {
		// 		@env()
		// 		name: "bin-\(_short)"
		// 		path: ("./bins/\(_bin)")
		// 		source: env.#Container & {
		// 			from: ctr.builder
		// 			steps: [
		// 				env.Env & {GOOS: _g, GOARCH: _a},
		// 				env.Exec & {args: ["go", "build", "-o", "./bins/\(_bin)", "./cmd/hof"]},
		// 			]
		// 		}
		// 	}
		// }
		{
			"\(_g)-\(_a)": env.#File & {
				@env()
				name: "bin-\(_g)-\(_a)"
				path: ("./bins/hof-\(_g)-\(_a)")
				source: env.#Container & {
					from: ctr.builder
					steps: [
						env.EnvVar & {GOOS: _g, GOARCH: _a},
						env.Exec & {args: ["go", "build", "-o", "./bins/hof-\(_g)-\(_a)", "./cmd/hof"]},
					]
				}
			}
		}
	}
}
