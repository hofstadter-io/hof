package workflows

import "github.com/hofstadter-io/hof/.github/workflows/common"

common.#Workflow & {
	name: "test_flow"
	on: pull_request: { paths: ["flow/**"] }
	jobs: test: {
		environment: "hof flow testing"
		steps: [ for step in common.#BuildSteps {step} ] + [{
			name: "Run mod tests"
			run: """
			hof flow -f test/flow
			"""
		}]
	}
}
