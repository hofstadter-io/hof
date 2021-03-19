package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	CompletionVanillaFlag bool
)

func init() {
	CompletionCmd.Flags().BoolVarP(&CompletionVanillaFlag, "vanilla", "8", false, "set to only check for an update")
}

var CompletionCmd = &cobra.Command{
	Use:     "completion",
	Aliases: []string{"completions"},
	Short:   "Generate completion helpers for popular terminals",
	Long:    "Generate completion helpers for popular terminals",
}

var BashCompletionLong = `Generate Bash completions

To load completion run

. <(hof completion bash)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(hof completion bash)
`

var BashCompletionCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate Bash completions",
	Long:  BashCompletionLong,
	Run: func(cmd *cobra.Command, args []string) {
		RootCmd.GenBashCompletion(os.Stdout)
	},
}

var ZshCompletionCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate Zsh completions",
	Long:  "Generate Zsh completions",
	Run: func(cmd *cobra.Command, args []string) {
		RootCmd.GenZshCompletion(os.Stdout)
	},
}

var FishCompletionCmd = &cobra.Command{
	Use:   "fish",
	Short: "Generate Fish completions",
	Long:  "Generate Fish completions",

	Run: func(cmd *cobra.Command, args []string) {
		RootCmd.GenFishCompletion(os.Stdout, true)
	},
}

var PowerShellCompletionCmd = &cobra.Command{
	Use:     "power-shell",
	Aliases: []string{"windows", "win", "power", "ps"},
	Short:   "Generate PowerShell completions",
	Long:    "Generate PowerShell completions",

	Run: func(cmd *cobra.Command, args []string) {
		RootCmd.GenPowerShellCompletion(os.Stdout)
	},
}

func init() {
	CompletionCmd.AddCommand(BashCompletionCmd)
	CompletionCmd.AddCommand(ZshCompletionCmd)
	CompletionCmd.AddCommand(FishCompletionCmd)
	CompletionCmd.AddCommand(PowerShellCompletionCmd)

	help := CompletionCmd.HelpFunc()
	usage := CompletionCmd.UsageFunc()

	CompletionCmd.SetHelpFunc(help)
	CompletionCmd.SetUsageFunc(usage)

}
