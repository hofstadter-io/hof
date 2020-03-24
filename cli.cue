package cli

import (
	"github.com/hofstadter-io/cuemod--cli-golang:cli"
	"github.com/hofstadter-io/cuemod--cli-golang/schema"
  "github.com/hofstadter-io/mvs"
)

Outdir: "./output/"

GEN: cli.Generator & {
	Cli: CLI
}

CLI: cli.Schema & {
	Name:    "hof"
	Package: "github.com/hofstadter-io/hof"

	Usage: "hof"
	Short: "hof is the cli for hof-lang, a low-code framework for developers"
	Long:  Short

  Releases: schema.GoReleaser & {
    Author: "Tony Worm"
    Homepage: "https://github.com/hofstadter-io/hof"

    Brew: {
      GitHubOwner: "hofstadter-io"
      GitHubRepoName: "homebrew-tap"
      GitHubUsername: "verdverm"
      GitHubEmail: "tony@hofstadter.io"
    }
  }

	OmitRun: true

	Pflags: [
		{
			Name:    "config"
			Type:    "string"
			Default: ""
			Help:    "Some config file path"
			Long:    "config"
			Short:   "c"
		},
		{
			Name:    "identity"
			Long:    "identity"
			Short:   "I"
			Type:    "string"
			Default: ""
			Help:    "the Studios Auth Identity to use during this hof execution"
		},
		{
			Name:    "context"
			Long:    "context"
			Short:   "C"
			Type:    "string"
			Default: ""
			Help:    "the Studios Context to use during this hof execution"
		},
		{
			Name:    "account"
			Long:    "account"
			Short:   "A"
			Type:    "string"
			Default: ""
			Help:    "the Studios Account to use during this hof execution"
		},
		{
			Name:    "project"
			Long:    "project"
			Short:   "P"
			Type:    "string"
			Default: ""
			Help:    "the Studios Project to use during this hof execution"
		},
	]

	PersistentPrerun: true
	PersistentPrerunBody: """
    fmt.Println("PersistentPrerun", RootConfigPflag, args)
  """

	Commands: [

		schema.Command & {
			Name:    "auth"
			Usage:   "auth"
			Short:   "authentication subcommands"
			Long:    Short
			OmitRun: true

			Commands: [
				schema.Command & {
					Name:  "login"
					Usage: "login"
					Short: "login to an account"
					Long:  Short

					Body: """
            fmt.Println("login not implemented")
          """
				},
			]
		},

		schema.Command & {
			Name:  "config"
			Usage: "config"
			Short: "configuration subcommands"
			Long:  Short

			OmitRun: true

			Commands: [
				schema.Command & {
					Name:  "list"
					Usage: "list"
					Short: "list configurations"
					Long:  Short
				},
				schema.Command & {
					Name:  "get"
					Usage: "get"
					Short: "print a configuration"
					Long:  Short
					Args: [
						schema.Arg & {
							Name:     "name"
							Type:     "string"
							Required: true
							Help:     "name of the configuration"
						},
					]
				},
				schema.Command & {
					Name:  "set"
					Usage: "set <name> <host> <account> [project]"
					Short: "set configuration values"
					Long:  Short
					Args: [
						schema.Arg & {
							Name:     "name"
							Type:     "string"
							Required: true
							Help:     "name for the configuration"
						},
						schema.Arg & {
							Name:     "host"
							Type:     "string"
							Required: true
							Help:     "host for this configuration"
						},
						schema.Arg & {
							Name:     "account"
							Type:     "string"
							Required: true
							Help:     "account for this configuration"
						},
						schema.Arg & {
							Name: "project"
							Type: "string"
							Help: "default project for this configuration"
						},
					]
				},
				schema.Command & {
					Name:  "use"
					Usage: "use"
					Short: "set the default configuration"
					Long:  Short
					Args: [
						schema.Arg & {
							Name:     "name"
							Type:     "string"
							Required: true
							Help:     "name of the configuration"
						},
					]
				},
			]
		},
		schema.Command & {
			Name:  "new"
			Usage: "new"
			Short: "create a new project or subcomponent or files"
			Long:  Short + ", depending on the context"
		},
		schema.Command & {
			Name:  "mod"
			Usage: "mod"
			Aliases: ["m"]
			Short: "manage project modules"
			Long:  "Hof has mvs embedded, so you can do all the same things from this subcommand"
      Commands: mvs.CLI.Commands
		},
		schema.Command & {
			Name:  "gen"
			Usage: "gen [files...]"
			Aliases: ["g"]
			Short: "generate code, data, and config"
			Long: """
        generate all the things, from code to data to config...
      """
		},
		schema.Command & {
			Name:  "studios"
			Usage: "studios"
			Aliases: ["s"]
			Short: "commands for working with Hofstadter Studios"
			Long: """
        Hofstadter Studios makes it easy to develop and launch both
        hof-lang modules as well as pretty much any code or application
      """

			OmitRun: true
			Commands: [

				schema.Command & {

					Name:  "secret"
					Usage: "secret"
					Aliases: [
						"secrets",
						"shh",
					]
					Short: "Work with Hofstadter Studios secrets"
					Long:  "Work with Hofstadter Studios secrets"

					Commands: [
						schema.Command & {
							Name:  "list"
							Usage: "list"
							Short: "List your secrets"
							Long:  "List your Studios secrets"
						},
						schema.Command & {
							Name:  "get"
							Usage: "get <name or id>"
							Short: "Get a Studios secret"
							Long:  Short
							Args: [ identArg]
						},
						schema.Command & {
							Name:  "create"
							Usage: "create <name> <input>"
							Short: "Create a Studios secret"
							Long:  "Create a Studios secret from a file or key/val pairs"
							Args: [
								nameArg,
								inputArg & {
									Help: "@file or key=val,key2=val2,..."
								},
							]
						},
						schema.Command & {
							Name:  "update"
							Usage: "update <name> <input>"
							Short: "Update a Studios secret"
							Long:  "Update a Studios secret from a file or key/val pairs"
							Args: [
								nameArg,
								inputArg & {
									Help: "@file or key=val,key2=val2,..."
								},
							]
						},
						schema.Command & {
							Name:  "delete"
							Usage: "delete <name or id>"
							Short: "Delete a Studios secret. Must not be in use"
							Long:  Short
							Args: [ identArg]
						},
					]
				},

			]

		},
		schema.Command & {
			Name:  "cue"
			Usage: "cue"
			Aliases: ["c"]
			Short: "Call a cue command"
			Long:  "Hof has cuelang embedded, so you can use hof anywhere you use cue"
			Body: """
        fmt.Println("run: cue", args)
      """
		},
	]
}

// Name arg
nameArg: schema.Arg & {
	Name:     "name"
	Type:     "string"
	Required: true
	Help:     "A name from /[a-zA-Z][a-zA-Z0-9_]*"
}

// Identifyier (name or id)
identArg: schema.Arg & {
	Name:     "ident"
	Type:     "string"
	Required: true
	Help:     "A name or id"
}

// input arg
inputArg: schema.Arg & {
	Name:     "input"
	Type:     "string"
	Required: true
}

contextArg: schema.Arg & {
	Name: "context"
	Type: "string"
	Help: "The hof auth context name"
}

// email for user / service account
identityArg: schema.Arg & {
	Name: "identity"
	Type: "string"
	Help: "A Hofstadter Studios user or service account"
}

// Studios account
accountArg: schema.Arg & {
	Name: "account"
	Type: "string"
	Help: "The name or id of a Hofstadter Studios account"
}

// Studios Project
projectArg: schema.Arg & {
	Name: "project"
	Type: "string"
	Help: "The name or id of a Hofstadter Studios project"
}

// Studios API Key
apikeyArg: schema.Arg & {
	Name: "apikey"
	Type: "string"
	Help: "Hofstadter Studios API Key"
}
