@experiment(aliasv2)
package veg

import (
	"strings"

	// "github.com/hofstadter-io/hof/lib/env/common/bases/lang"
	"github.com/hofstadter-io/hof/schemas/env"
)

flags: {
  // TODO, scope these since we are now injecting the entire package at the repo root
	repo: string | *"https://github.com/hofstadter-io/hof" @tag(repo)
	// todo, change this to "." when we move something to the index, if we ever really do?
	local:  string | *"."              @tag(local)
	src:    "repo" | *"local" | string @tag(src,short=repo|local)
	adk:    string | *"../adk"         @tag(adk)
	dagger: string | *"../dagger"      @tag(dagger)

	goos: string | *"darwin" @tag(goos,var=os)
	arch: string | *"arm64" @tag(arch,var=arch)

	use: {
		lsp: bool | *false
	}

	ports: {
		gopls:  int | *4000 @tag(ports_gopls)
		cuepls: int | *4001 @tag(ports_cuepls)
	}
}

src: {
	repo: env.#GitRepo & {
		@env()
		name: string | *"repo"
		url:  flags.repo
	}
	local: env.#HostDir & {
		@env()
		name: string | *"local"
		path: flags.local
	}
	adk: env.#HostDir & {@env(), path: flags.adk}
	dagger: env.#HostDir & {@env(), path: flags.dagger}

	// setup code base on flags and value
	code: {@env(), name: string}
	if flags.src == "repo" {code: repo}
	if flags.src == "local" {code: local}
	if strings.HasPrefix(flags.src, "https://") {
		code: env.#GitRepo & {url: flags.src}
	}
	if strings.HasPrefix(flags.src, ".") {
		code: env.#HostDir & {path: flags.src}
	}

}

incept: {
	[string]~(K,_): {name: K}
	docker: env.#HostSocket & {@env(), path: "unix:///Users/tony/.colima/default/docker.sock"}
	dagger: env.#HostSocket & {@env(), path: "unix:///var/run/docker.sock"}
}

out: {
	// this should filter from either repo or host
	cuemod: env.#Dir & {
		@env()
		name:   "cuemod"
		path:   "."
		source: src.repo
		include: [
			"cue.mod",
			"schemas",
			"examples",
			"lib/env/common",
			// way more to come here
		]
	}

	// go builds, cross-arch/os
	cli: {}

	vscode: {}

	docs: {}

	fmtrs: {}
}
