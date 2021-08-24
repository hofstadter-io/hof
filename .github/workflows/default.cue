import "github.com/hofstadter-io/hof/.github/workflows/common"

common.#Workflow & {
	name: "hof-test"
	on: ["pull_request"]
	jobs: test: {
		steps: [ for step in common.#BuildSteps {step} ] + [{
			name: "Run gen tests"
			run: """
			hof test test.cue -s gen -t test
			"""
		},{
			name: "Run tester tests"
			run: """
			hof test test/testers/api/postman.cue
			"""
		}]
	}
}
