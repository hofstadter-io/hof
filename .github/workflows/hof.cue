package workflows

import "github.com/hofstadter-io/ghacue"

#Workflow: ghacue.#Workflow & {
	name: string
	on:   _ | *["push", "pull_request", "workflow_dispatch"]
	env: HOF_TELEMETRY_DISABLED: "1"
	jobs: test: {
		strategy: matrix: {
			"go-version": ["1.18.x", "1.19.x"]
			os: ["ubuntu-latest", "macos-latest"]
		}
		strategy: "fail-fast": false
		"runs-on": "${{ matrix.os }}"
	}
}

#Workflow & {
	name: "hof"
	jobs: test: {
		environment: "hof mod testing"
		steps:
			[ for step in #BuildSteps {step} ] +
			[ for step in #TestSteps  {step} ]
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
},{
	name: "Set up Docker"
	uses: "docker/setup-buildx-action@v2"
}, {
	name: "Fetch Go deps"
	run:  "go mod download"
}, {
	name: "Build CLI"
	run:  "go install ./cmd/hof"
}, {
	name: "Build Formatters"
	run:  "make formatters"
}]

#TestSteps: [{
	name: "Run self gen test"
	run: """
		# fetch CUE deps
		hof mod vendor cue
		# generate templates
		hof gen
		# should have no diff
		git diff
		# git diff --exit-code
		"""
	env: {
		HOFMOD_SSHKEY:      "${{secrets.HOFMOD_SSHKEY}}"
		GITHUB_TOKEN:       "${{secrets.HOFMOD_TOKEN}}"
	}
}, {
	// maybe these should be services?
	name: "Start formatting containers"
	run:  "hof fmt start"
}, {
	name: "Run template test"
	run: """
		hof flow @test/gen ./test.cue
		"""
}, {
	name: "Run render tests"
	run: """
		hof flow @test/render ./test.cue
		"""
	env: {
		HOFMOD_SSHKEY:      "${{secrets.HOFMOD_SSHKEY}}"
		GITHUB_TOKEN:       "${{secrets.HOFMOD_TOKEN}}"
	}
}, {
	name: "Run lib/structural tests"
	run: """
		hof flow @test/st ./test.cue
		"""
}, {
	name: "Run flow tests"
	run: """
		hof flow -f test/flow ./test.cue
		"""
}, {
	name: "Run mod tests"
	run: """
		hof flow -f test/mod ./test.cue
		"""
	env: {
		HOFMOD_SSHKEY:      "${{secrets.HOFMOD_SSHKEY}}"
		GITHUB_TOKEN:       "${{secrets.HOFMOD_TOKEN}}"
		GITLAB_TOKEN:       "${{secrets.GITLAB_TOKEN}}"
		BITBUCKET_USERNAME: "hofstadter"
		BITBUCKET_PASSWORD: "${{secrets.BITBUCKET_TOKEN}}"
	}
}, {
	// should probably be last?
	name: "Run fmt tests"
	run: """
		hof flow -f test/fmt ./test.cue
	"""
}]
