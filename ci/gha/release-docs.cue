package workflows

import (
	"github.com/hofstadter-io/ghacue"
	"github.com/hofstadter-io/hof/ci/gha/common"
)

ghacue.#Workflow & {
	name: "release (docs)"

	on: {
		push: {
			tags: ["docs-**"]
		}
		workflow_dispatch: {
			inputs: {
				deploy: {
					description: "where to deploy"
					required:    true
					default:     "next"
					type:        "choice"
					options: ["next", "prod"]
				}
			}
		}
	}
	env: HOF_TELEMETRY_DISABLED: "1"

	jobs: {
		docs: {
			"runs-on":   "ubuntu-latest"
			environment: "hof docs"

			steps: [
				// general setup
				common.Steps.cue.install,
				common.Steps.go.setup & {#ver: "1.20.x"},
				common.Steps.go.cache,
				common.Steps.checkout,
				common.Steps.vars,
				common.Steps.go.deps,
				common.Steps.hof.install,

				// prod build site & image
				common.Steps.docs.setup,
				common.Steps.docs.env,
				{
					name: "Build"

					run: """
						cd docs
						make gen
						make hugo.${DOCS_ENV}
						"""
				},

				// gcloud auth setup
				common.Steps.gcloud.auth,
				common.Steps.gcloud.setup,
				common.Steps.gcloud.dockerAuth,

				// push image & deploy
				{
					name: "Image"
					run: """
						export TAG=${HOF_TAG}

						cd docs
						make docker
						make push
						make deploy.${DOCS_ENV}.view
						"""
					// todo, we need to create a CloudBuild and run that
				},
			]
		}
	}
}
