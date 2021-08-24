import "github.com/hofstadter-io/hof/.github/workflows/common"

common.#Workflow & {
	name: "test_mod"
	on: pull_request: {
		paths: ["lib/mod/**"]
	}
	jobs: test: {
		steps: [ for step in common.#BuildSteps {step} ] + [{
			name: "Run mod tests"
			run: """
			hof test test.cue -s lib -t test -t mod
			"""
			env: {
				GITHUB_TOKEN: "${{secrets.HOFMOD_TOKEN}}"
			}
		}]
	}
}

