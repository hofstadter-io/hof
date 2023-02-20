package workflows

import (
	"github.com/hofstadter-io/ghacue"
	"github.com/hofstadter-io/hof/ci/gha/common"
)

ghacue.#Workflow & {
	name: "hof"
	on:   _ | *["push", "pull_request", "workflow_dispatch"]
	env: HOF_TELEMETRY_DISABLED: "1"
	jobs: test: {
		strategy: {
			"fail-fast": false
			matrix: {
				"go": [...] & common.Versions.go
				os: [...] & common.Versions.os
			}
		}
		environment: "hof mod testing"
		"runs-on": "${{ matrix.os }}"

		steps: [
			common.Steps.go.setup & { #ver: "${{ matrix.go }}" },
			common.Steps.go.cache,
			common.Steps.checkout,
			common.Steps.vars,
			common.Steps.go.deps,
			{
				name: "Build CLI"
				run:  "go install ./cmd/hof"
			},
			common.Steps.docker.setup & {
				"if": "${{ !startsWith( runner.os, 'macos') }}"
			},
			common.Steps.docker.compat & {
				"if": "${{ !startsWith( runner.os, 'macos') }}"
			},
			{
				name: "Build Formatters"
				run:  """
					make formatters
					hof fmt info
					"""
				"if": "${{ !startsWith( runner.os, 'macos') }}"
			},
		] + #TestSteps
	}
}

#TestSteps: [...{
	env: {
		HOFMOD_SSHKEY:      "${{secrets.HOFMOD_SSHKEY}}"
		GITHUB_TOKEN:       "${{secrets.HOFMOD_TOKEN}}"
	}
}]

#TestSteps: [{
	name: "test/self"
	run: """
		# fetch CUE deps
		hof mod vendor cue
		# generate templates
		hof gen
		# should have no diff
		git diff
		# git diff --exit-code
		"""
}, {
	// maybe these should be services?
	name: "Start formatters"
	run:  """
		hof fmt start
		hof fmt info
		"""
	"if": "${{ !startsWith( runner.os, 'macos') }}"
}, {
	name: "test/gen"
	run: """
		hof flow @test/gen ./test.cue
		"""
}, {
	name: "test/render"
	run: """
		hof flow @test/render ./test.cue
		"""
}, {
	name: "test/structural"
	run: """
		hof flow @test/st ./test.cue
		"""
}, {
	name: "test/flow"
	run: """
		hof flow -f test/flow ./test.cue
		"""
}, {
	name: "test/mod"
	run: """
		hof flow -f test/mod ./test.cue
		"""
	env: {
		GITLAB_TOKEN:       "${{secrets.GITLAB_TOKEN}}"
		BITBUCKET_USERNAME: "hofstadter"
		BITBUCKET_PASSWORD: "${{secrets.BITBUCKET_TOKEN}}"
	}
}, {
	// should probably be last?
	name: "test/fmt"
	run: """
		hof fmt info
		hof flow -f test/fmt ./test.cue
		"""
	"if": "${{ !startsWith( runner.os, 'macos') }}"
}]
