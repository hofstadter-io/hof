package tool

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

agents: {
	lsp2mcp: env.Exec & {
		args: ["sh", "-c", _script]

		_script: """
			# LSP -> MCP
			go install github.com/isaacphi/mcp-language-server@latest
			"""
	}
}
