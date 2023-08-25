package workflows

import (
	"github.com/hofstadter-io/ghacue"
	"github.com/hofstadter-io/hof/ci/gha/common"
)

ghacue.#Workflow & {
	name: "hof"
	on:   _ | *["push"]
	env: HOF_TELEMETRY_DISABLED: "1"
	jobs: test: {
		concurrency: {
			group:                "${{ github.workflow }}-${{ matrix.os }}-${{ matrix.go }}-${{ github.ref_name }}"
			"cancel-in-progress": true
		}
		strategy: {
			"fail-fast": false
			matrix: {
				"go": [...] & common.Versions.go
				os:   [...] & common.Versions.os
			}
		}
		environment: "hof mod testing"
		"runs-on":   "${{ matrix.os }}"

		steps: [
			common.Steps.go.setup & {#ver: "${{ matrix.go }}"},
			common.Steps.go.cache,
			common.Steps.checkout,
			common.Steps.vars,

			// common.Steps.docker.setup,
			common.Steps.docker.machack,
			common.Steps.docker.macos,
			common.Steps.docker.compat,

			// application steps
			common.Steps.go.deps,
			{
				name: "Build CLI"
				run:  "go install ./cmd/hof"
			},
			{
				name: "Build Formatters"
				run: """
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
		GITHUB_TOKEN: "${{secrets.HOFMOD_TOKEN}}"
	}
}]

#TestSteps: [{
	name: "test/self"
	run: """
		# self: gen -> diff
		set -e

		# mods & deps
		hof mod tidy
		hof fmt cue.mod/module.cue
		hof mod vendor

		# gen self
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
	name: "test/datamodel"
	run: """
		hof flow @test/dm ./test.cue
		"""
}, {
	name: "test/flow"
	run: """
		hof flow -f test/flow ./test.cue
		"""
}, {
	name: "test/fmt"
	run: """
		docker ps -a
		hof fmt info
		hof flow -f test/fmt ./test.cue
		"""
}, {
	// should probably be last for external workflows?
	// or maybe separate workflow for permissions?
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
}]
