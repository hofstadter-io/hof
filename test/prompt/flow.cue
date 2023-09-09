package prompt

import "encoding/json"

test: {
	@flow()

	prompt: {
		@task(prompt.Prompt)

		Input: {
			name:     string
			repo:     string
			releases: bool | *false
		}

		Questions: [{
			Name:     "name"
			Type:     "input"
			Prompt:   "What is your project named"
			Required: true
		}, {
			Name:    "repo"
			Type:    "input"
			Prompt:  "Git repository"
			Default: "github.com/user/repo"
		}, {
			Name:    "releases"
			Type:    "confirm"
			Prompt:  "Enable release tooling"
			Default: true
		}]

		Output: _
	}

	output: {
		@task(os.Stdout)
		text: json.Indent(json.Marshal(prompt.Output), "", "  ")
	}

}
