package workflows

import "github.com/hofstadter-io/hof/.github/workflows/common"

common.#Workflow & {
	name: "default"
	on: ["push"]
	jobs: test: {
		steps: [ for step in common.#BuildSteps {step}] + [{
			name: "Run structural tests"
			run: """
				hof flow @test/gen ./test.cue
				hof flow @test/render ./test.cue
				hof flow @test/st ./test.cue
				"""
		}]
	}
}
