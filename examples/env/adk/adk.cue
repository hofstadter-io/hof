@experiment(aliasv2)
package adk

import (
	"github.com/hofstadter-io/hof/lib/env/common/bases/lang"
	"github.com/hofstadter-io/hof/schemas/env"
)

_flags: {
	repo: string | *"https://github.com/google/adk-go" @tag(repo)
	// eventually this will go at the root of the repo
	code: string | *"/Users/tony/adk/go" @tag(code)
	// app: string | *"/Users/tony/adk/go" @tag(app)
}

src: {
	[string]~(k,_): {@env(), name: k}
	repo: env.#Dir & {path: ".", source: env.#GitRepo & {url: _flags.repo}}
	code: env.#HostDir & {path: _flags.code}
	// app: env.#HostDir & { path: _flags.app }

	// _actual: repo
}

ctr: {
	[string]~(k,_): {@env(), name: k}
	base: env.#Container & {
		from: lang.go.ctr.base
		steps: [
			env.Mount & {path: "/work", source: src.repo},
		]
	}
	dev: env.#Container & {
		from: base
		steps: [
			env.BindService & {alias: "gopls", service: lang.go.svc.gopls & {name: "gopls", source: from}},
		]
	}
}

_tester: env.#Container & {
	#cmd: string
	from: ctr.base
	steps: [
		env.Exec & {args: ["bash", "-c", _script]},
	]
	_script: """
  set -euo pipefail
  \(#cmd)
  """
}

cmd: {
	[string]~(k1,_): env.#Cmd & {
		@env(), name: k1
		tasks: [string]~(k2,_): {
			@env(), name: k2
			steps: [[{name: "\(k1).\(k2)"}]]
		}
	}

	test: tasks: {
		go: {steps: [[_tester & {#cmd: "go test ./..."}]]}
		race: {steps: [[_tester & {#cmd: "go test -race ./..."}]]}
		cover: {steps: [[_tester & {#cmd: "go test -cover ./..."}]]}
	}
	lint: tasks: {
		// want something like: gofmt -l . | wc -l | grep -e '^0$'
		fmt: {steps: [[_tester & {#cmd: #"gofmt -l . || true"#}]]}
		staticcheck: {steps: [[_tester & {#cmd: "staticcheck ./... || true"}]]}
		golangci: {steps: [[_tester & {#cmd: "golangci-lint run || true"}]]}
	}
}
