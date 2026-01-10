package tool

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

github: {
	cli: env.Sh & {
		// https://github.com/cli/cli/blob/trunk/docs/install_linux.md#debian
		script: """
			mkdir -p -m 755 /etc/apt/keyrings
			out=$(mktemp)
			wget -nv -O$out https://cli.github.com/packages/githubcli-archive-keyring.gpg
			cat $out | tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null 
			chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg
			mkdir -p -m 755 /etc/apt/sources.list.d
			echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null

			apt-get update -y
			apt-get install -y gh
			"""
	}
}
