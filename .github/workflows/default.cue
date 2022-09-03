package workflows

import "github.com/hofstadter-io/hof/.github/workflows/common"

common.#Workflow & {
	name: "default"
	on: ["push", "workflow_dispatch"]
	jobs: test: {
		steps: [ for step in common.#BuildSteps {step}] + [{
			name: "Run self-gen test"
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
