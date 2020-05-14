package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	// "github.com/spf13/viper"

	"github.com/hofstadter-io/mvs/lib"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/cmd/hof/pflags"
)

var hofLong = `Polyglot Code Gereration Framework`

func init() {

	RootCmd.PersistentFlags().StringVarP(&pflags.RootConfigPflag, "config", "c", "", "Some config file path")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootIdentityPflag, "identity", "I", "", "the Studios Auth Identity to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootContextPflag, "context", "C", "", "the Studios Context to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootAccountPflag, "account", "A", "", "the Studios Account to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootProjectPflag, "project", "P", "", "the Studios Project to use during this hof execution")
}

func RootPersistentPreRun(args []string) (err error) {

	lib.InitLangs()

	return err
}

func RootPersistentPostRun(args []string) (err error) {

	PrintUpdateAvailable()

	return err
}

var RootCmd = &cobra.Command{

	Use: "hof",

	Short: "Polyglot Code Gereration Framework",

	Long: hofLong,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = RootPersistentPreRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendGaEvent("root", "<omit>", 0)

	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = RootPersistentPostRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	hf := RootCmd.HelpFunc()
	f := func(cmd *cobra.Command, args []string) {
		if RootCmd.Name() == cmd.Name() {
			ga.SendGaEvent("root/help", "<omit>", 0)
		}
		hf(cmd, args)
	}
	RootCmd.SetHelpFunc(f)

	cobra.OnInitialize(initConfig)
	RootCmd.AddCommand(AuthCmd)
	RootCmd.AddCommand(ConfigCmd)
	RootCmd.AddCommand(NewCmd)
	RootCmd.AddCommand(GenCmd)
	RootCmd.AddCommand(StudiosCmd)
	RootCmd.AddCommand(ModCmd)
	RootCmd.AddCommand(RunCmd)
	RootCmd.AddCommand(CueCmd)
}

func initConfig() {

}
