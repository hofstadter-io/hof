package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

ModCmdImports: [
	{Path: "github.com/hofstadter-io/hof/lib/mod", ...},
	{Path: "github.com/hofstadter-io/hof/cmd/hof/flags", ...},
]

ModCommand: schema.Command & {
	// TBD:   "β"
	Name:  "mod"
	Usage: "mod"
	Aliases: ["m"]
	Short: "CUE module dependency management"
	Long:  "CUE module dependency management, imported from upstream CUE project"

	//Topics: #ModTopics
	//Examples: #ModExamples
}
