package tool

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

cloud: {
	gcloud: env.Exec & {
		args: ["sh", "-c", _script]

		// https://docs.cloud.google.com/sdk/docs/install-sdk#deb
		_script: """
			set -eou pipefail

			curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg
			echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
			apt-get update -y
			apt-get install -y google-cloud-cli
			"""
	}

	awscli: env.Exec & {
		args: ["sh", "-c", _script]

		_arch: "x86_64" | "aarch64"

		// https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
		_script: """
			set -eou pipefail

			curl "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip"
			unzip awscliv2.zip
			./aws/install
			rm -rf awscliv2.zip ./aws
			"""
	}

	azure: env.Exec & {
		args: ["sh", "-c", _script]

		// https://learn.microsoft.com/en-us/cli/azure/install-azure-cli-linux?view=azure-cli-latest&pivots=apt
		_script: """
			set -eou pipefail

			curl -sL https://aka.ms/InstallAzureCLIDeb | bash
			"""
	}
}
