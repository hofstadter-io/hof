package workflows

import "github.com/hofstadter-io/hof/.github/workflows/common"

common.#Workflow & {
	name: "default"
	jobs: test: {
		steps: [ for step in common.#BuildSteps {step}] + [{
			name: "Run self gen test"
			run: """
				# fetch CUE deps
				hof mod vendor cue
				# generate templates
				hof gen
				# should have no diff
				git diff --exit-code
				"""
		},{
			name: "Run template test"
			run: """
				hof flow @test/gen ./test.cue
				"""
		},{
			name: "Run render tests"
			run: """
				hof flow @test/render ./test.cue
				"""
		},{
			name: "Run lib/structural tests"
			run: """
				hof flow @test/st ./test.cue
				"""
		}]
	}
}
