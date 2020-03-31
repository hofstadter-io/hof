package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "local"
	Commit  = "dirty"

	BuildDate = "unknown"
	GoVersion = "unknown"
	BuildOS   = "unknown"
	BuildArch = "unknown"
	BuildArm  = "n/a"
)

const versionMessage = `
Version:     v%s
Commit:      %s

BuildDate:   %s
GoVersion:   %s
OS / Arch:   %s %s
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
		fmt.Printf(
			versionMessage,
			Version,
			Commit,
			BuildDate,
			GoVersion,
			BuildOS,
			BuildArch,
			BuildArm,
		)
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}
