package hof

import (
	"github.com/hofstadter-io/cuemod--cli-golang:cli"
	"github.com/hofstadter-io/cuemod--cli-golang/schema"
)

Outdir :: "./"

GenCli: cli.Generator & {
	Cli: CLI
}

_LibImport :: [
	schema.Import & {Path: CLI.Package + "/lib"},
]

CLI :: cli.Schema & {
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


	Imports: [
		schema.Import & {Path: "github.com/hofstadter-io/mvs/lib"},
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

