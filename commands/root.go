package commands

import (
	"github.com/spf13/cobra"
)

var HofLong = `HofLang is a language and transpiler
for building data-centric DSLs and designs.
`

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
