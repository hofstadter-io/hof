name: "hof cli"

image: "mcr.microsoft.com/devcontainers/universal:2"

postCreateCommand: """
make hof
"""

customizations: {
	vscode: extensions: [
		"asdine.cue",
		"jallen7usa.vscode-cue-fmt",
	]
}
