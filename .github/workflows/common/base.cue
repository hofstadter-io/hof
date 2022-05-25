package common

import "github.com/hofstadter-io/ghacue"

#Workflow: ghacue.#Workflow & {
	name: string
	on:   _ | *["pull_request"]
	jobs: test: {
		strategy: matrix: {
			"go-version": ["1.17.x", "1.18.x"]
			os: ["ubuntu-latest", "macos-latest"]
		}
		strategy: "fail-fast": false
		"runs-on": "${{ matrix.os }}"
	}
}

#BuildSteps: [{
	name: "Install Go"
	uses: "actions/setup-go@v2"
	with: "go-version": "${{ matrix.go-version }}"
}, {
	name: "Checkout code"
	uses: "actions/checkout@v2"
}, {
	name: "Build CLI"
	run:  "go install ./cmd/hof"
}]
