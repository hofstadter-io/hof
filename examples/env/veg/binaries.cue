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
			env.Dir & {path: "/adk", source: src.adk.fork},
			env.Dir & {path: "/dagger", source: src.dagger.fork},
			env.Dir & {path: "/work", source: env.#HostDir & {
				path: "."
				name: "hof-bin-src"
				include: [
					"go.mod", "go.sum", "cue.mod",
					"cmd", "lib", "flow", "script",
				]
			}},
		]
	}
}

// File version
hof: File: { for k,v in hof.matrix { (k): env.File & { path: "/usr/local/bin/hof", content: v}}}
// #File verions
hof: #File: {
	[string]~(k,_): {name: "bin-\(k)"}

	_maker: env.#File & {
		#variant: string
		#goos: string
		#arch: string
		@env()
		#hof: id: "hof-cli-\(#variant)"

		source: env.#Container & {
			from: ctr.builder
			steps: [
				env.EnvVar & {GOOS: #goos, GOARCH: #arch},
				env.Exec & {args: ["go", "build", "-o", "./bins/hof", "./cmd/hof"]},
			]
		}
	}

	local: _maker & {#goos: flags.goos, #arch: flags.arch}
	matrix: {
		_goos: ["linux", "darwin"]
		_arch: ["amd64", "arm64"]
		for _g in _goos {
			// local, by goos
			(_g): _maker & {#variant: _g, #goos: _g, #arch: flags.arch}
			// arch x goos
			for _a in _arch {
				"\(_g)-\(_a)": _maker & {#variant: "\(_g)-\(_a)", #goos: _g, #arch: _a}
			}
		}
	}
}
