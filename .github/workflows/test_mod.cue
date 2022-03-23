package workflows

import "github.com/hofstadter-io/hof/.github/workflows/common"

common.#Workflow & {
	name: "test_mod"
	on: pull_request: { paths: ["lib/mod/**", "lib/yagu/repos/**", "lib/yagu/git.go", "lib/yagu/netrc.go", "lib/yagu/ssh.go"] }
	jobs: test: {
		environment: "hof mod testing"
		steps: [ for step in common.#BuildSteps {step} ] + [{
			name: "Run mod tests"
			run: """
			hof flow -f test/mods ./test.cue
			"""
			env: {
				HOFMOD_SSHKEY: "${{secrets.HOFMOD_SSHKEY}}"
				GITHUB_TOKEN: "${{secrets.HOFMOD_TOKEN}}"
				GITLAB_TOKEN: "${{secrets.GITLAB_TOKEN}}"
				BITBUCKET_USERNAME: "hofstadter"
				BITBUCKET_PASSWORD: "${{secrets.BITBUCKET_TOKEN}}"
			}
		}]
	}
}

