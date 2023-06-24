package workflows

import (
	"github.com/hofstadter-io/ghacue"
	"github.com/hofstadter-io/hof/ci/gha/common"
)

ghacue.#Workflow & {
	name: "hack"
	on:   _ | *["push", "pull_request", "workflow_dispatch"]
	env: {
		HOF_TELEMETRY_DISABLED: "1"
		HOF_FMT_VERSION:        "v0.6.8-rc.5"
	}
	jobs: {
		dagger: {
			environment: "hof mod testing"
			"runs-on":   "ubuntu-latest"

			steps: [
				common.Steps.go.setup & {#ver: "1.20.x"},
				common.Steps.go.cache,
				common.Steps.checkout,
				common.Steps.vars,
				common.Steps.go.deps,
				common.Steps.docker.setup,
				common.Steps.docker.macos,
				common.Steps.docker.compat,

				{
					name: "dagger-in-dagger"
					run:  "go run ./test/dagger/main/dagger-in-dagger.go"
				},
				{
					name: "dockerd-in-dagger"
					run:  "go run ./test/dagger/main/dockerd-in-dagger.go"
				},
				{
					name: "dagger-in-dagger"
					run:  "go run ./test/dagger/main/hof.go"
					env: {
						GITHUB_TOKEN: "${{secrets.HOFMOD_TOKEN}}"
					}
				},
			]
		}
	}
}
