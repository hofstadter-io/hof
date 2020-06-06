package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
	"github.com/hofstadter-io/hof/cmd/hof/verinfo"
)

const versionMessage = `
Version:     v%s
Commit:      %s

BuildDate:   %s
GoVersion:   %s
OS / Arch:   %s %s


Author:   Hofstadter, Inc
Homepage: https://hofstadter.io
GitHub:   https://github.com/hofstadter-io/hof

`

var VersionLong = `Print the build version for hof`

var VersionCmd = &cobra.Command{

	Use: "version",

	Aliases: []string{
		"ver",
	},

	Short: "print the version",

	Long: VersionLong,

	Run: func(cmd *cobra.Command, args []string) {

		s, e := os.UserConfigDir()
		fmt.Printf("hof ConfigDir %q %v\n", filepath.Join(s, "hof"), e)

		fmt.Printf(
			versionMessage,
			verinfo.Version,
			verinfo.Commit,
			verinfo.BuildDate,
			verinfo.GoVersion,
			verinfo.BuildOS,
			verinfo.BuildArch,
		)
	},
}

func init() {
	help := VersionCmd.HelpFunc()
	usage := VersionCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		if VersionCmd.Name() == cmd.Name() {
			ga.SendCommandPath("version help")
		}
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		if VersionCmd.Name() == cmd.Name() {
			ga.SendCommandPath("version usage")
		}
		return usage(cmd)
	}
	VersionCmd.SetHelpFunc(thelp)
	VersionCmd.SetUsageFunc(tusage)

}
