@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

_tester: env.#Container & {
	#cmd: string
	from: env.#Container & {
		@id(tester-with-src)
		from: root.ctr.dev
		steps: [
			env.Dir & {path: "/adk", source: src.adk},
			env.Dir & {path: "/dagger", source: src.dagger},
			env.Dir & {path: "/work", source: src.code},
		]
	}
	steps: [
		env.Bash & {script: "\(#cmd)"},
	]
}

cmd: {
	[string]~(k1,_): env.#Cmd & {
		@env(), name: k1
		#hof: id: "cmd-\(k1)"
		#hof: metadata: name: #hof.id

		tasks: [string]~(k2,_): {
			name: k2
			steps: [...[...{name: "\(k1).\(k2)"}]]
		}
	}

	// most of this should move to the Go pack
	test: tasks: {
		go: steps: [[_tester & {#cmd: "go test ./..."}]]
		govet: steps: [[_tester & {#cmd: "go vet ./..."}]]
		goveti: steps: [[_tester & {#cmd: "go test ./..."}]]
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
	// 	full:    steps: [test, lint, scan, review]
	// 	release: steps: [full, gather, publish]

	// 	// env.#HostExec (todo)
	// 	gather: ["hof env export -P dist"]
	// 	publish: [
	// 		"git tag",
	// 		"gh cli to draft & upload",
	// 	]

	// 	onPush: default
	// 	prPush: full
	// 	onTag:  steps: [release]

	// }
}
