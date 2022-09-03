package workflows

import "github.com/hofstadter-io/hof/.github/workflows/common"

common.#Workflow & {
	name: "test_flow"
	_paths: [
		".github/workflows/test_flow.*",
		"flow/**",
	]
	on: {
		workflow_dispatch: {}
		pull_request: {paths: _paths}
		push: {paths: _paths}
	}
	jobs: test: {
		environment: "hof flow testing"
		steps:       [ for step in common.#BuildSteps {step}] + [{
			name: "Run mod tests"
			run: """
				hof flow -f test/flow ./test.cue
				"""
		}]
	}
}
