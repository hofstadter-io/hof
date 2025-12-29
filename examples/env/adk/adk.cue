@experiment(aliasv2)
package adk

import (
	"github.com/hofstadter-io/hof/lib/env/common/bases/lang"
	"github.com/hofstadter-io/hof/schemas/env"
)

_flags: {
	repo: string | *"https://github.com/google/adk-go" @tag(repo)
	// eventually this will go at the root of the repo
	code: string | *"." @tag(code)
}

src: {
	[string]~(k,_): {@env(), name: k}
	repo: env.#GitRepo & { url: _flags.repo }
	code: env.#HostDir & { path: _flags.code }
}

ctr: {
	[string]~(k,_): {@env(), name: k}
	base: env.#Container & {
		from: lang.go.ctr.base
		steps: [
			env.Mount & {path: "/work", source: src.code},
		]
	}
	dev: env.#Container & {
		from: base
		steps: [
			env.BindService & { alias: "gopls", service: lang.go.svc.gopls & { name: "gopls", source: from }},
		]
	}

}

cmd: {
	[string]~(k,_): env.#Cmd & {@env(), name: k}
	test: {
		steps: [[]]
	}
}
