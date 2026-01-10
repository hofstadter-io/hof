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

	astralSteps: [
		env.Sh & {
			script: """
				# uv 
				curl -LsSf https://astral.sh/uv/install.sh | env UV_INSTALL_DIR="/usr/local/bin" sh

				uv tool install ruff@latest
				uv tool install ty@latest
				"""
		},
	]
	defaultSteps: [
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
	defaultExtras: [
		env.Sh & {
			// yes, pyright requires node and recommends installing via npm
			script: """
				# LSP
				npm install -g pyright
				"""
		},
	]
}
