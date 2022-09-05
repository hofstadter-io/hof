package common

import "github.com/hofstadter-io/ghacue"

#Workflow: ghacue.#Workflow & {
	name: string
	on:   _ | *["push", "pull_request", "workflow_dispatch"]
	jobs: test: {
		strategy: matrix: {
			"go-version": ["1.18.x", "1.19.x"]
			os: ["ubuntu-latest", "macos-latest"]
		}
		strategy: "fail-fast": false
		"runs-on": "${{ matrix.os }}"
	}
}

#BuildSteps: [{
	name: "Install Go"
	uses: "actions/setup-go@v3"
	with: "go-version": "${{ matrix.go-version }}"
}, {
	uses: "actions/cache@v3"
	with: {
		path: #"""
			~/go/pkg/mod
			~/.cache/go-build
			~/Library/Caches/go-build
			~\AppData\Local\go-build
			"""#
		key:            "${{ runner.os }}-go-${{ matrix.go-version }}-${{ hashFiles('**/go.sum') }}"
		"restore-keys": "${{ runner.os }}-go-${{ matrix.go-version }}-"
	}
}, {
	name: "Checkout code"
	uses: "actions/checkout@v3"
}, {
	name: "Fetch Go deps"
	run:  "go mod download"
}, {
	name: "Build CLI"
	run:  "go install ./cmd/hof"
}]
