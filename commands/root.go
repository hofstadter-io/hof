package commands

import (
	"github.com/spf13/cobra"

	"cuelang.org/go/cue/load"
)

var HofLong = `HofLang is a language and transpiler
for building data-centric DSLs and designs.
`

var (
	CueConfig *load.Config
)

func init() {
	CueConfig = &load.Config{}
}

var (
	RootCmd = &cobra.Command{

		Use: "hof",

		Short: "HofLang framework CLI tool",

		Long: HofLong,

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Argument Parsing

		},
	}
)
