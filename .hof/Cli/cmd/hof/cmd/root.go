package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	// "github.com/spf13/viper"

	"github.com/hofstadter-io/mvs/lib"

	"github.com/hofstadter-io/hof/lib/runtime"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/cmd/hof/pflags"
)

var hofLong = `Polyglot Code Gereration Framework`

func init() {

	RootCmd.PersistentFlags().StringVarP(&pflags.RootConfigPflag, "config", "C", "", "Path to a hof configuration file")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootAccountPflag, "account", "A", "", "the account context to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootBillingPflag, "billing", "B", "", "the billing context to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootProjectPflag, "project", "P", "", "the project context to use during this hof execution")
	RootCmd.PersistentFlags().StringSliceVarP(&pflags.RootLabelsPflag, "label", "L", nil, "Labels for use across all commands")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootVerbosityPflag, "verbosity", "v", "", "set the verbosity of output")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootQuietPflag, "quiet", "q", false, "turn off output and assume defaults at prompts")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootImpersonateAccountPflag, "impersonate-account", "I", "", "account to impersonate for this hof execution, relies on having permission to impersonate and avoids need for credentials")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootTraceTokenPflag, "trace-token", "", "", "used to help debug issues")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootLogHTTPPflag, "log-http", "", "", "used to help debug issues")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootCredsPflag, "creds", "", "", "The path to a hof creds file")
}

func RootPersistentPreRun(args []string) (err error) {

	lib.InitLangs()
	runtime.Init()

	return err
}

func RootPersistentPostRun(args []string) (err error) {

	WaitPrintUpdateAvailable()

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
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(AuthCmd)
	RootCmd.AddCommand(ConfigCmd)
	RootCmd.AddCommand(NewCmd)
	RootCmd.AddCommand(GenCmd)
	RootCmd.AddCommand(StudiosCmd)
	RootCmd.AddCommand(ModCmd)
	RootCmd.AddCommand(RunCmd)
	RootCmd.AddCommand(CueCmd)
	RootCmd.AddCommand(ModelsetCmd)
	RootCmd.AddCommand(StoreCmd)
	RootCmd.AddCommand(ImportCmd)
	RootCmd.AddCommand(ExportCmd)
	RootCmd.AddCommand(UiCmd)
	RootCmd.AddCommand(HackCmd)
}

func initConfig() {

}
