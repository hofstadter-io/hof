package workflows

import "github.com/hofstadter-io/ghacue"
import "github.com/hofstadter-io/hof/.github/workflows/common"

ghacue.#Workflow & {
	name: "windows"
	on: ["push"]
	jobs: test: {
		strategy: matrix: {
			"go-version": ["1.16.x", "1.17.x"]
			os: ["windows-latest"]
		}
		strategy: "fail-fast": false
		"runs-on": "${{ matrix.os }}"
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
		},{
			name: "Run mod tests"
			run: """
			hof test test.cue -s lib -t test -t mod
			"""
			env: {
				GITHUB_TOKEN: "${{secrets.HOFMOD_TOKEN}}"
				GITLAB_TOKEN: "${{secrets.GITLAB_TOKEN}}"
			}
		}]
	}
}

