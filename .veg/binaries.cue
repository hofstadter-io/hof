package veg

import (
	"github.com/hofstadter-io/hof/catalogs/env/packs/lang"
	"github.com/hofstadter-io/hof/schemas/env"
)

// go multi-target builds
// from source build too for local setup

ctr: {
	builder: env.#Container & {
		@env(hof-cli-builder)
		from: lang.go.ctr.base
		steps: [
			env.Dir & {path: "/adk", source: src.adk.fork},
			env.Dir & {path: "/dagger", source: src.dagger.fork},
			env.Dir & {path: "/work", source: env.#HostDir & {
				path: flags.disk
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
hof: File: {for k, v in hof.cli.matrix {
	(k): env.File & {
		#hof: {id: v.#hof.id, metadata: v.#hof.metadata}
		name:    v.name
		path:    string | *"/usr/local/bin/hof"
		content: v
	}
}}
// #File verions
hof: cli: {
	_maker: env.#File & {
		// metadata
		@env()
		#hof: id: "hof-cli-\(#variant)"
		#hof: metadata: name: #hof.id

		// params
		#variant: string
		#goos:    string
		#arch:    string

		path: "./bins/hof-\(#variant)"
		source: env.#Container & {
			from: ctr.builder
			steps: [
				env.EnvVars & {GOOS: #goos, GOARCH: #arch},
				env.Exec & {args: ["go", "build", "-ldflags", "-w", "-o", "./bins/hof-\(#variant)", "./cmd/hof"]},
			]
		}
	}

	local: _maker & {#variant: "local", #goos: flags.goos, #arch: flags.arch}
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
