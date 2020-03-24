package commands

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"
)

var hofLong = `hof is the cli for hof-lang, a low-code framework for developers`

var (
	RootConfigPflag   string
	RootIdentityPflag string
	RootContextPflag  string
	RootAccountPflag  string
	RootProjectPflag  string
)

func init() {

	RootCmd.PersistentFlags().StringVarP(&RootConfigPflag, "config", "c", "", "Some config file path")

	RootCmd.PersistentFlags().StringVarP(&RootIdentityPflag, "identity", "I", "", "the Studios Auth Identity to use during this hof execution")

	RootCmd.PersistentFlags().StringVarP(&RootContextPflag, "context", "C", "", "the Studios Context to use during this hof execution")

	RootCmd.PersistentFlags().StringVarP(&RootAccountPflag, "account", "A", "", "the Studios Account to use during this hof execution")

	RootCmd.PersistentFlags().StringVarP(&RootProjectPflag, "project", "P", "", "the Studios Project to use during this hof execution")

}

var RootCmd = &cobra.Command{

	Use: "hof",

	Short: "hof is the cli for hof-lang, a low-code framework for developers",

	Long: hofLong,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		lib.InitLangs()

	},
}

func init() {
	RootCmd.AddCommand(AuthCmd)
	RootCmd.AddCommand(ConfigCmd)
	RootCmd.AddCommand(NewCmd)
	RootCmd.AddCommand(ModCmd)
	RootCmd.AddCommand(GenCmd)
	RootCmd.AddCommand(StudiosCmd)
	RootCmd.AddCommand(CueCmd)
}
