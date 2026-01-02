package lang

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

node: {
	#ver:  string | *"24.12.0"
	#arch: *"arm64" | "amd64"

	caches: {
		nodeMods: env.Volume & {
			name: "node-mods-\(#ver)-\(#arch)"
			type: "cache"
		}
	}

	default: [
		// cache
		install,
	]

	install: [
		env.Exec & {

			_src:  "https://nodejs.org/dist/v\(#ver)/\(_file)"
			_file: "node-v\(#ver)-linux-\(#arch).tar.xz"

			_script: """
        set -eou pipefail

        cd /tmp
        wget -q \(_src)
        tar -C /usr/local -xf \(_file) --strip-components=1
        rm -rf /tmp/*

        # LSP
        npm install -g yarn pnpm tsx typescript typescript-language-server
        """

			args: ["sh", "-c", _script]

		},
	]
}
