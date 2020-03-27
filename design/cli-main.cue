package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

Outdir :: "./"

_LibImport :: [
	schema.Import & {Path: CLI.Package + "/lib"},
]

CLI :: schema.Cli & {
	Name:    "hof"
	Package: "github.com/hofstadter-io/hof"

	Usage: "hof"
	Short: "hof is the cli for hof-lang, a low-code framework for developers"
	Long:  Short

	Releases: schema.GoReleaser & {
		Author:   "Tony Worm"
		Homepage: "https://github.com/hofstadter-io/hof"

		Brew: {
			GitHubOwner:    "hofstadter-io"
			GitHubRepoName: "homebrew-tap"
			GitHubUsername: "verdverm"
			GitHubEmail:    "tony@hofstadter.io"
		}
	}

	OmitRun: true

  Pflags: CliPflags

  EnablePProf: true


	Imports: [
		{Path: "github.com/hofstadter-io/mvs/lib"},
	]

	PersistentPrerun: true
	PersistentPrerunBody: """
    lib.InitLangs()
  """

	Commands: [
    // meta
    AuthCommand,
    ConfigCommand,

    // hof
    NewCommand,
    GenCommand,
    StudiosCommand,

    // extern
    CmdCommand,
    ModCommand,
    CueCommand,

	]
}

