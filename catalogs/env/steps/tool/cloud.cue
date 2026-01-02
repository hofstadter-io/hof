package tool

import (
	"github.com/hofstadter-io/hof/catalogs/env/bases"
	"github.com/hofstadter-io/hof/schemas/env"
)

cloud: {
	// apt/debian
	gcloud: env.Sh & {
		// https://docs.cloud.google.com/sdk/docs/install-sdk#deb
		script: """
			set -eou pipefail

			curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg
			echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
			apt-get update -y && apt-get install -y google-cloud-cli
			"""
	}

	awscli: aws.fetch

	aws: {
		#arch: *"aarch64" | "x86_64"
		fetch: env.Sh & {

			// https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
			script: """
				curl "https://awscli.amazonaws.com/awscli-exe-linux-\(#arch).zip" -o "awscliv2.zip"
				unzip awscliv2.zip
				./aws/install
				rm -rf awscliv2.zip ./aws
				"""
		}
		out: env.#Dir & {
			sources: [
				env.#Container & {
					from: bases.debian.default
					steps: [fetch]
				},
			]
			include: [
				"/usr/local/aws-cli",
				"/usr/local/bin/aws",
				"/usr/local/bin/aws_completer",
			]
		}
		dir: env.Dir & {
			source: out
		}

	}

	// https://learn.microsoft.com/en-us/cli/azure/install-azure-cli-linux?view=azure-cli-latest&pivots=apt
	azure: env.Sh & {script: "curl -sL https://aka.ms/InstallAzureCLIDeb | bash"}
}
