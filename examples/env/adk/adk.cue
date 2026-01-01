@experiment(aliasv2)
package adk

import (
	"github.com/hofstadter-io/hof/lib/env/common/bases/lang"
	"github.com/hofstadter-io/hof/schemas/env"
)

_flags: {
	// eventually this will go at the root of the repo and just be "."
	local: string | *"/Users/tony/hof/adk"               @tag(local)
	fork:  string | *"https://github.com/verdverm/adk-go" @tag(fork)
	repo:  string | *"https://github.com/google/adk-go" @tag(repo)

	// are we using source from local or git
	use: "repo" | *"local" @tag(use,short=repo|local)
}

src: {
	[string]~(k,_): {@env(), name: k}
	repo: env.#Dir & {path: ".", sources: [env.#GitRepo & {url: _flags.repo}]}
	local: env.#HostDir & {path: _flags.local}
	// app: env.#HostDir & { path: _flags.app }

	_actual: _
	if _flags.use == "local" {_actual: local}
	if _flags.use == "repo" {_actual: repo}
}

ctr: {
	[string]~(k,_): {@env(), name: k}
	base: env.#Container & {
		from: lang.go.ctr.base
		steps: [
			env.Mount & {path: "/work", source: src._actual},
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
			steps: [...[...{name: "\(k1).\(k2)"}]]
		}
	}

	test: tasks: {
		go: steps: [[_tester & {#cmd: "go test ./..."}]]
		// parallel tests
		goUltra: steps: [[
			_tester & {#cmd: "go vet ./..."},
			_tester & {#cmd: "go test -race ./..."},
			_tester & {#cmd: "go test -cover ./..."},
		]]
		// sequential tests
		// vet: {steps: [[_tester & {#cmd: "go vet ./..."}]]}
		// race: {steps: [[_tester & {#cmd: "go test -race ./..."}]]}
		// cover: {steps: [[_tester & {#cmd: "go test -cover ./..."}]]}
	}
	lint: tasks: {
		// want something like: gofmt -l . | wc -l | grep -e '^0$'
		fmt: steps: [[_tester & {#cmd: #"gofmt -l . || true"#}]]
		staticcheck: steps: [[_tester & {#cmd: "staticcheck ./... || true"}]]
		golangci: steps: [[_tester & {#cmd: "golangci-lint run || true"}]]
	}
}
