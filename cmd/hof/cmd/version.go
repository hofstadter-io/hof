package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
	"github.com/hofstadter-io/hof/cmd/hof/verinfo"
	"github.com/hofstadter-io/hof/lib/container"
)

const versionMessage = `hof - the high code framework

Version:     %s
Commit:      %s

BuildDate:   %s
GoVersion:   %s
CueVersion:  %s
OS / Arch:   %s %s
ConfigDir:   %s
CacheDir:    %s
Containers:  %s

Author:      Hofstadter, Inc
Homepage:    https://hofstadter.io
GitHub:      https://github.com/hofstadter-io/hof
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

		configDir, _ := os.UserConfigDir()
		cacheDir, _ := os.UserCacheDir()

		err := container.InitClient()
		if err != nil {
			fmt.Println(err)
		}

		rt, err := container.GetVersion()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf(
			versionMessage,
			verinfo.Version,
			verinfo.Commit,
			verinfo.BuildDate,
			verinfo.GoVersion,
			verinfo.CueVersion,
			verinfo.BuildOS,
			verinfo.BuildArch,
			filepath.Join(configDir,"hof"),
			filepath.Join(cacheDir,"hof"),
			rt,
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
