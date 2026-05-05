package hof

import (
	"github.com/hofstadter-io/hof/.veg:veg"
)

veg

flags: {
	// TODO, scope these since we are now injecting the entire package at the repo root
	repo: string | *"https://github.com/hofstadter-io/hof" @tag(repo)

	// todo, change this to "." when we move something to the index, if we ever really do?
	disk:   string | *"."             @tag(disk)
	src:    "repo" | *"disk" | string @tag(src,short=repo|disk)
	ref:    string | *"_next"         @tag(ref)
	adk:    string | *"../adk"        @tag(adk)
	dagger: string | *"../dagger"     @tag(dagger)

	goos: string @tag(goos,var=os)
	arch: string @tag(arch,var=arch)

	lsp: bool | *false

	ports: {
		gopls:  int | *4000 @tag(ports_gopls)
		cuepls: int | *4001 @tag(ports_cuepls)
	}

	registry: string | *"host.docker.internal:5000" @tag(reg)
}
