@experiment(aliasv2)
package basic

import "github.com/hofstadter-io/hof/schemas/env"

multi: {
	#ver:  string | *"1.25" @tag(ver)
	#arch: string           @tag(arch,var=arch)
	#os:   string           @tag(os,var=os)

	// base container for building in
	base: env.#Container & {
		@env(multi-base)
		from: "golang:\(#ver)-alpine"
		steps: [
			// mount caches for mods and intermediate build artifacts (saves time)
			env.Mount & {path: "/cache/go", source: env.#Cache & {name: "go-build-\(#ver)-\(#arch)"}},
			env.Mount & {path: "/go", source: env.#Cache & {name: "go-mods-\(#ver)-\(#arch)"}},

			// set any default Go vars for all builds
			env.EnvVars & {CGO_ENABLED: "0"},

			// a globally consistent workdir
			env.Workdir & {path: "/work"},
		]
	}

	// a container after the code has built
	built: env.#Container & {
		@env(multi-built)
		from: base
		steps: [
			// here we are passing the content directly as a string
			// load and add directories from your filesystem with #HostDir
			env.File & {path: "main.go", content: _goSrc},

			// run the build
			env.Exec & {args: ["go", "build", "-ldflags", "-w -s", "-o", "server", "main.go"]},
		]
	}

	// actual binary file for the server
	binary: env.#File & {@env(multi-binary), path: "server", source: built}

	// start from an alpine, add the binary
	runner: env.#Container & {@env(multi-runner)
		from: "alpine:latest"
		steps: [
			// Sh is a wrapper around Exec to `sh -c <script>`
			env.Sh & {script: "apk add --update --no-cache ca-certificates"},

			// add the server binary
			env.File & {path: "/usr/bin/server", content: binary},

			// set the entrypoint
			env.Entrypoint & {args: ["/usr/bin/server"]},
		]
	}

	service: env.#Service & {
		@env(multi-service)
		ports: [{port: 8080}]
		source: runner
	}

	// caches we mount to the Go toolchain for across session caching
	caches: {
		// [string]~(k,_): {@env(), name: "go-\(k)"}
		build: env.Mount & {path: "/cache/go", source: env.#Cache & {name: "go-build-\(#ver)-\(#arch)"}}
		mods: env.Mount & {path: "/go", source: env.#Cache & {name: "go-mods-\(#ver)-\(#arch)"}}
	}
}

// you can @embed() files in CUE or load them with veg/env using #HostDir
_goSrc: """
	package main
	
	import (
	    "fmt"
	    "net/http"
	)
	
	func main() {
	    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	        fmt.Fprintf(w, "Hello, you've requested: %s\\n", r.URL.Path)
	    })
	
	    http.ListenAndServe(":8080", nil)
	}
	"""
