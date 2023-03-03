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
			common.Steps.docker.setup,
			common.Steps.docker.macos,
			common.Steps.docker.login,
			common.Steps.docker.compat,
			{
				name: "Build Formatters"
				run:  """
					make formatters
					docker images
					hof fmt start
					hof fmt info
					docker ps -a
					"""
			},
		] + #TestSteps
	}
}

#TestSteps: [...{
	env: {
		GITHUB_TOKEN:       "${{secrets.HOFMOD_TOKEN}}"
	}
}]

#TestSteps: [{
	name: "test/self"
	run: """
		# self: gen -> diff
		set -e

		# mods stuff & gen
		hof mod tidy
		hof mod vendor
		hof gen

		# should have no diff
		git diff --exit-code
		"""
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
	name: "test/create"
	run: """
		hof flow @test/create ./test.cue
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
		HOFMOD_SSHKEY:      "${{secrets.HOFMOD_SSHKEY}}"
		GITLAB_TOKEN:       "${{secrets.GITLAB_TOKEN}}"
		BITBUCKET_USERNAME: "hofstadter"
		BITBUCKET_PASSWORD: "${{secrets.BITBUCKET_TOKEN}}"
	}
}, {
	// should probably be last?
	name: "test/fmt"
	run: """
		docker ps -a
		hof fmt info
		hof flow -f test/fmt ./test.cue
		"""
	env: {
		HOF_FMT_DEBUG: "1"
	}
}]
