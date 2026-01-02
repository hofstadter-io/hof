package lang

import (
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

python: {
	caches: {
		pythonMods: env.Volume & {
			name: "python-mods"
			type: "cache"
		}
	}

	default: [
		utils.apt.install & {#pkgs: [
			"pip",
			"pipx",
			"pipenv",
			"pylint",
			"python3-poetry",
			"python3-pytest",
			"python3-flake8",
		]},
	]
	dev: [
		env.Exec & {
			args: ["sh", "-c", _script]

			// yes, pyright requires node and recommends installing via npm
			_script: """
				set -eou pipefail

				# uv 
				curl -LsSf https://astral.sh/uv/install.sh | env UV_INSTALL_DIR="/usr/local/bin" sh

				# LSP
				npm install -g pyright
				"""
		},
	]
}
