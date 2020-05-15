package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#Outdir: "./cmd/hof"
#Module: "github.com/hofstadter-io/hof"

#_LibImport: [
	schema.#Import & {Path: #Module + "/lib"},
]

#CLI: schema.#Cli & {
	Name:    "hof"
	Package: "github.com/hofstadter-io/hof/cmd/hof"

	Usage: "hof"
	Short: "Polyglot Code Gereration Framework"
	Long:  Short

	Releases: #CliReleases

  Telemetry: "UA-103579574-5"
  TelemetryIdDir: "hof"

	OmitRun: true

  Pflags: #CliPflags

  PersistentPrerun: true
  PersistentPostrun: true
  // EnablePProf: true


	Imports: [
		{Path: "github.com/hofstadter-io/mvs/lib"},
		{Path: "github.com/hofstadter-io/hof/lib/runtime"},
	]

	PersistentPrerun: true
	PersistentPrerunBody: """
    lib.InitLangs()
		runtime.Init()
  """

	Commands: [
    // meta
    #InitCommand,
    #AuthCommand,
    #ConfigCommand,

    // hof
    #NewCommand,
    #GenCommand,
    #StudiosCommand,

    // extern
    #ModCommand,
    #RunCommand,
    #CueCommand,

    #ModelCommand,
    #StoreCommand,

    #ImportCommand,
    #ExportCommand,

    #UiCommand,

    // for dev
    schema.#Command & {
      Name:    "hack"
      Usage:   "hack ..."
      Short:   "development command"
      Long: Short
    },
	]
}

