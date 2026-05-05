@experiment(aliasv2)
package adk

import (
	"github.com/hofstadter-io/hof/catalogs/env/packs"
	"github.com/hofstadter-io/hof/schemas/env"
)

flags: {
	goos:   string @tag(goos,var=os)
	goarch: string @tag(goarch,var=arch)
}

_goPack: packs.lang.go & {#goos: flags.goos, #goarch: flags.goarch}

src: {
	[string]~(k,_): {@env(), name: k}
	main: env.#GitRepo & {url: "https://github.com/google/adk-go", ref: "main"}
	fork: env.#GitRepo & {url: "https://github.com/verdverm/adk-go", ref: "veg"}
	code: env.#Dir & {path: ".", sources: [fork]}
}

ctr: {
	[string]~(k,_): {@env(), name: k}
	base: env.#Container & {
		from: packs.lang.go.ctr.base
		steps: [
			env.Mount & {path: "/work", source: src.code},
		]
	}
	dev: env.#Container & {
		from: base
		steps: [
			env.BindService & {alias: "gopls", service: packs.lang.go.svc.lsp},
		]
	}
}

_tester: env.#Container & {
	#cmd: string
	from: ctr.base
	steps: [env.Bash & {script: "\(#cmd)"}]
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
		spelling: _
	}

	scan: tasks: {
		sonar: {}
		vuln: {}
	}

	review: tasks: {
		agent: {
			// ... code changes,
			// docs / agents.md need updating,
			// stage & apply suggested changes,
		}
	}

	// ci: tasks: {
	// 	default: steps: [test, lint]
	// 	full: steps: [test, lint, scan, review]
	// 	release: steps: [full, gather, publish]

	// 	// env.#HostExec (todo)
	// 	gather: ["hof env export -P dist"]
	// 	publish: [
	// 		"git tag",
	// 		"gh cli to draft & upload",
	// 	]

	// 	onPush: default
	// 	prPush: full
	// 	onTag: steps: [release]

	// }
}
