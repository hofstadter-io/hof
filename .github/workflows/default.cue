import "github.com/hofstadter-io/ghacue"

ghacue.#Workflow & {
	name: "hof"
	on: ["push", "pull_request"]
	jobs: test: {
		strategy: matrix: {
			"go-version": ["1.15.x", "1.16.x"]
			os: ["ubuntu-latest", "macos-latest", "windows-latest"]
		}
		"runs-on": "${{ matrix.os }}"
		steps: [{
			name: "Install Go"
			uses: "actions/setup-go@v2"
			with: "go-version": "${{ matrix.go-version }}"
		},{
			name: "Checkout code"
			uses: "actions/checkout@v2"
		},{
			name: "Build hof"
			run: """
			go mod download
			GOOS=${{ matrix.os }} go build ./cmd/hof
			"""
		}]
	}
}
