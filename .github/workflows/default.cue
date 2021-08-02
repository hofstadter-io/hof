import "github.com/hofstadter-io/ghacue"

ghacue.#Workflow & {
	name: "hof"
	on: ["push", "pull_request"]
	jobs: test: {
		strategy: matrix: {
			"go-version": ["1.15.x", "1.16.x"]
			os: ["ubuntu-latest", "macos-latest", "windows-latest"]
		}
		strategy: "fail-fast": false
		"runs-on": "${{ matrix.os }}"
		steps: [{
			name: "Install Go"
			uses: "actions/setup-go@v2"
			with: "go-version": "${{ matrix.go-version }}"
		},{
			name: "Checkout code"
			uses: "actions/checkout@v2"
		},{
			name: "Download mods"
			run: "go mod download"
		},{
			name: "Build CLI"
			run: "go install ./cmd/hof"
		},{
			name: "Run tests"
			run: "hof test test.cue"
		}]
	}
}
