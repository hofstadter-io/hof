package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#ModCmdImports: [
	{Path: "github.com/hofstadter-io/hof/lib/mod", ...},
]

#ModCommand: schema.#Command & {
  Name:  "mod"
  Usage: "mod"
  Aliases: ["m"]
	Short: "mod subcmd is a polyglot dependency management tool based on go mods"
	Long: """
  The mod subcmd is a polyglot dependency management tool based on go mods.

  mod file format:

    module = "<module path>"

    <name> = "version"

    require (
      ...
    )

    replace <module path> => <local path>
    ...
  """

	OmitRun: true

	Imports: #ModCmdImports

	PersistentPrerun: true
	PersistentPrerunBody: """
    mod.InitLangs()
  """
	Commands: [
		{
			Name:  "info"
			Usage: "info [language]"
			Short: "print info about languages and modders known to mvs"
			Long: """
        print info about languages and modders known to mvs
          - no arg prints a list of known languages
          - an arg prints info about the language modder configuration that would be used
      """

			Args: [
				{
					Name: "lang"
					Type: "string"
					Help: "name of the language to print info about"
				},
			]

			Imports: #ModCmdImports

			Body: """
      msg, err := mod.LangInfo(lang)
      if err != nil {
        fmt.Println(err)
        os.Exit(1)
      }
      fmt.Println(msg)
      """
		},
		{
			Name:  "convert"
			Usage: "convert <lang> <file>"
			Short: "convert another package system to MVS."
			Long:  Short

			Args: [
				{
					Name:     "lang"
					Type:     "string"
					Required: true
					Help:     "name of the language to print info about"
				},
				{
					Name:     "filename"
					Type:     "string"
					Required: true
					Help:     "existing package filename, depending on language"
				},
			]

			Imports: #ModCmdImports

			Body: """
      err = mod.Convert(lang, filename)
      if err != nil {
        fmt.Println(err)
        os.Exit(1)
      }
      """
		},
		{
			Name:  "graph"
			Usage: "graph"
			Short: "print module requirement graph"
			Long:  Short

			Imports: #ModCmdImports

			Body: """
      err = mod.ProcessLangs("graph", args)
      if err != nil {
        fmt.Println(err)
        os.Exit(1)
      }
      """
		},
		{
			Name:  "status"
			Usage: "status"
			Short: "print module dependencies status"
			Long:  Short

			Imports: #ModCmdImports

			Body: """
      err = mod.ProcessLangs("status", args)
      if err != nil {
        fmt.Println(err)
        os.Exit(1)
      }
      """
		},
		{
			Name:  "init"
			Usage: "init <lang> <module>"
			Short: "initialize a new module in the current directory"
			Long:  Short

			Args: [
				{
					Name:     "lang"
					Type:     "string"
					Required: true
					Help:     "name of the language to print info about"
				},
				{
					Name:     "module"
					Type:     "string"
					Required: true
					Help:     "module name or path, depending on language"
				},
			]

			Imports: #ModCmdImports

			Body: """
      err = mod.Init(lang, module)
      if err != nil {
        fmt.Println(err)
        os.Exit(1)
      }
      """
		},
		{
			Name:  "tidy"
			Usage: "tidy [langs...]"
			Short: "add missinad and remove unused modules"
			Long:  Short

			Imports: #ModCmdImports

			Body: """
      err = mod.ProcessLangs("tidy", args)
      if err != nil {
        fmt.Println(err)
        os.Exit(1)
      }
      """
		},
		{
			Name:  "vendor"
			Usage: "vendor [langs...]"
			Short: "make a vendored copy of dependencies"
			Long:  Short

			Imports: #ModCmdImports

			Body: """
      err = mod.ProcessLangs("vendor", args)
      if err != nil {
        fmt.Println(err)
        os.Exit(1)
      }
      """
		},
		{
			Name:  "verify"
			Usage: "verify [langs...]"
			Short: "verify dependencies have expected content"
			Long:  Short

			Imports: #ModCmdImports

			Body: """
      err = mod.ProcessLangs("verify", args)
      if err != nil {
        fmt.Println(err)
        os.Exit(1)
      }
      """
		},
	]

}
