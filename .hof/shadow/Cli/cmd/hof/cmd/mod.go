package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/mod"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/mod"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var modLong = `The mod subcmd is a polyglot dependency management tool based on go mods.

mod file format:

  module = "<module path>"

  <name> = "version"

  require (
    ...
  )

  replace <module path> => <local path>
  ...`

func ModPersistentPreRun(args []string) (err error) {

	mod.InitLangs()

	return err
}

var ModCmd = &cobra.Command{

	Use: "mod",

	Aliases: []string{
		"m",
	},

	Short: "mod subcmd is a polyglot dependency management tool based on go mods",

	Long: modLong,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ModPersistentPreRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},
}

func init() {

	help := ModCmd.HelpFunc()
	usage := ModCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	ModCmd.SetHelpFunc(thelp)
	ModCmd.SetUsageFunc(tusage)

	ModCmd.AddCommand(cmdmod.InfoCmd)
	ModCmd.AddCommand(cmdmod.ConvertCmd)
	ModCmd.AddCommand(cmdmod.GraphCmd)
	ModCmd.AddCommand(cmdmod.StatusCmd)
	ModCmd.AddCommand(cmdmod.InitCmd)
	ModCmd.AddCommand(cmdmod.TidyCmd)
	ModCmd.AddCommand(cmdmod.VendorCmd)
	ModCmd.AddCommand(cmdmod.VerifyCmd)

}

var ModTopics = map[string]string{
	"cue": `hof mod for Cuelang is a prototype until module and dependency management is in the Cue tool.

The version here is quite robust, however it lacks the code introspection like Golang,
so you are required to manage your own cue.mods file.
`,
	"custom": `hof mod allows you to create custom module systems and have MVS dependency management out of the box.

see: hof mod --example custom
`,
	"go": `hof mod for Golang uses the Go tool behind the scenes.
`,
	"js": `Not supported yet, but the idea is to use Go's MVS system on JS and other languages.
`,
	"mod-file": `The mod file is the same format as Golang mod files. The only difference is that
the "Go = 1.14" is replaced with "<lang/mod> = <major>.<minor>"
`,
	"py": `Not supported yet, but the idea is to use Go's MVS system on Python and other languages.
`,
	"readme": `# hof mod - a golang MVS system for anything

¡hof mod¡ is a flexible tool and library based on Go mods.

Use and create module systems with [Minimum Version Selection](https://research.swtch.com/vgo-mvs) semantics
and manage dependencies ¡go mod¡ style.
Mix any set of language, code bases, git repositories, package managers, and subdirectories.
Manage polyglot and monorepo codebase dependencies with
[100% reproducible builds](https://github.com/golang/go/wiki/Modules#version-selection) from a single tool.


### Features

- Based on go mods MVS system, aiming for 100% reproducible builds.
- Recursive dependencies, version resolution, and code instrospection.
- Custom module systems with custom file names and vendor directories.
- Control configuration for naming, vendoring, and other behaviors.
- Polyglot support for multiple module systems and multiple languages within one tool.
- Works with any git system and supports the main features from go mods.
- Convert other vendor and module systems to MVS or just manage their files with MVS.

Language support:

- [golang](https://golang.org) - delegates to the go tool for the commands mirrored here
- [cuelang](https://cuelang.org) - builtin in default using the custom module feature
- [hofmods](https://hofstadter.io) - extends Cue with advanced code generation capabilities
- [custom](./docs/custom-modders.md) - Create your own locally or globally with ¡.mvsconfig¡ files

Upcoming languages: Python and JavaScript
so they can have an MVS system and the benefits,
and ¡hof mod¡ can start supporting and fetching from package registries.
These language implementations will have flexibility to
manage with ¡hof mod¡ and the common toolchain to varying degrees.
Pull requests for improved language support are welcome.

The main difference from go mods is that ¡hof mod¡, generally,
is not introspecting your code to determine dependencies.
It relies on user management of the ¡<lang>.mods¡ file.
Since golang is exec'd out to, introspection is supported,
and as more languages improve, we look to similarly
improve this situation in ¡hof mod¡.
A midstep to full custom implementation will be a
introspection custom module with some basic support
using file globs and regex lists.

_Note, we also default to the plural ¡<lang>.mods/sums¡ files and ¡<lang.mod>/¡ vendor directory.
This is 1) because cuelang reads from ¡cue.mod¡ directory, and 2) because it is less likely
to collide with existing directories.
You can also configure more behavior per language and module than go mods.
The go mods is shelled out to as a convience, and often languages impose restrictions._


### Usage

¡¡¡shell
# Print known languages in the current directory
hof mod info

# Initialize this folder as a module
hof mod init <lang> <module-path>

# Add your requirements
vim <lang>.mods  # go.mod like file

# Pull in dependencies, no args discovers by *.mods and runs all
hof mod vendor [langs...]

# See all of the commands
hof mod help
¡¡¡


### Module File

The module file holds the requirements for project.
It has the same format as a ¡go.mod¡ file.

¡¡¡
# These are like golang import paths
#   i.e. github.com/hofstadter-io/hof
module <module-path> 

# Information about the module type / version
#  some systems that take this into account
# go = 1.14
<lang> = <version>

# Required dependencies section
require (
	# <module-path> <module-semver>
	github.com/hof-lang/cuemod--cli-golang v0.0.0      # This is latest on HEAD
	github.com/hof-lang/cuemod--cli-golang v0.1.5      # This is a tag v0.1.5 (can omit 'v' in tag, but not here)
)

# replace <module-path> => <module-path|local-path> [version]
replace github.com/hof-lang/cuemod--cli-golang => github.com/hofstadter-io/hofmod-cli-golang v0.2.0
replace github.com/hof-lang/cuemod--cli-golang => ../../cuelibs/cuemod--cli-golang
¡¡¡


### Custom Module Systems

¡.mvsconfig.cue¡ allows you to define custom module systems.
With some simple configuration, you can create and control
and vendored module system based on ¡go mods¡.
You can also define global configurations.

See the [custom-modder docs](./docs/custom-modders.md)
to learn more about writing
you own module systems.

This is the current Cue __modder__ configuration:

¡¡¡cue
cue: {
	Name: "cue"
	Version: "0.1.1"
	ModFile: "cue.mods"
	SumFile: "cue.sums"
	ModsDir: "cue.mod/pkg"
	MappingFile: "cue.mod/modules.txt"
	InitTemplates: {
		"cue.mod/module.cue": """
			module "{{ .Module }}"
			"""
		}

	VendorIncludeGlobs: []
	VendorExcludeGlobs: [
		"/.git/**",
		"**/cue.mod/pkg/**",
	]
}
¡¡¡

### Motivation

- MVS has better semantics for vendoring and gets us closer to 100% reproducible builds.
- JS and Python can have MVS while still using the remainder of the tool chains.
- Prototype for cuelang module and vendor management.
- We need a module system for our [hof-lang](https://hof-lang.org) project.

#### Links about go mods

[Using go modules](https://blog.golang.org/using-go-modules)

[Go and Versioning](https://research.swtch.com/vgo)

[More about version selection](https://github.com/golang/go/wiki/Modules#version-selection)


#### Other

You may also like to look at the [hofmod-cli](https://github.com/hofstadter-io/hofmod-cli) project.
We use this to generate the CLI code and files for CI.

You can find us in the
[cuelang slack](https://join.slack.com/t/cuelang/shared_invite/enQtNzQwODc3NzYzNTA0LTAxNWQwZGU2YWFiOWFiOWQ4MjVjNGQ2ZTNlMmIxODc4MDVjMDg5YmIyOTMyMjQ2MTkzMTU5ZjA1OGE0OGE1NmE)
for now.

"""

"overview": ##"""
hof mods enable you to manage dependencies and versions for just about anything.

You can run ¡hof mod <op> <lang> ...¡ for language specific operations or use
just ¡hof mod <op>¡ to auto-discover module systems and run the operation on each.

Modules do not have to be tied to a language, they can be arbitrary git repos.
You have a lot of control for which files to include/exclude as well as
life-cycle operations for performing any tasks around your modules.
`,
}

var ModExamples = map[string]string{
	"cue": `# Create a Cue module:

A Cue module example.

### Create a Cue module:

¡¡¡
hof mod init cue github.com/hofstadter-io/my-cue-mod
¡¡¡

### Vendor dependencies

¡¡¡
hof mod vendor cue
¡¡¡

You can use replaces to work on code locally
`,
	"custom": `# Custom Modules and Dependency Systems

hod mod gives you the ability to create
custom module systems, called Modders.
Modder is the struct name for the
internal code which controls how
modules and vendoring is handled.
You can configure as many of these as you like,
by providing global or local ¡.mvsconfig.cue¡ files.

You create your own modders by createing ¡.mvsconfig.cue¡ files.
hog mod will look for these in two places before any commands are run.

- A global ¡$HOME/.mvs/config.cue¡
- A local ¡./.mvsconfig.cue¡

### hof mod "modder" config file format


¡¡¡cue
// These two need to be the same
cue: {
	Name: "cue"
	// non-semver of the language
	Version: "#.#.#"

	// Common defaults, can be anything
	ModFile:  "<lang>.mods"
	SumFile:  "<lang>.sums"
	ModsDir:  "<lang>.mod/pkg"
	Checksum: "<lang>.mod/checksum.txt"

	// Controls for modders that want to shell out
	// to common tools for certain commands
	NoLoad: false
	CommandInit: [[string]]
	CommandGraph: [[string]]
	CommandTidy: [[string]]
	CommandVendor: [[string]]
	CommandVerify: [[string]]
	CommandStatus: [[string]]

	// Runs on init for this language
	// filename/content key/pair values
	// uses the golang text/template library
	// inputs will be
	//   .Language
	//   .Module
	//   .Modder
	InitTemplates: {
		"<lang>.mod/module.<lang>": """
					module "{{ .Module }}"
					"""
	}
	// Series of commands to be executed pre/post init
	InitPreCommands: [[string]]
	InitPostCommands: [[string]]

	// Same as the InitTemplates, but run during vendor command
	VendorTemplates: {
		"<lang>.mod/module.<lang>": """
					module "{{ .Module }}"
					"""
	}

	VendorIncludeGlobs: [
		".mvsconfig.cue",
		"<lang>.mods",
		"<lang>.sums",
		"<lang>.mod/module.<lang>",
		"<lang>.mod/modules.txt",
		"**/*.<lang>",
	]
	VendorExcludeGlobs: ["<lang>.mod/pkg"]

	// Series of commands to be executed pre/post vendoring
	VendorPreCommands: [[string]]
	VendorPostCommands: [[string]]

	// Use MVS to only manage the languages normal dependency file
	ManageFileOnly: false

	// Whether local replaces should use a symlink instead of copying files
	SymlinkLocalReplaces: false

	// Controls the code introspection for dependency determiniation
	IntrospectIncludeGlobs: ["**/*.<lang>"]
	IntrospectExcludeGlobs: ["<lang>.mod/pkg"]
	IntrospectExtractRegex: ["you will have to figure out a series of 'any match passes' regexps to pull out dependencies"]

	// This field determines the prefix to place in front of
	// imports which have a single token or leverage package managers
	// This is currently futurology for building MVS for Python and JavaScript
	PackageManagerDefaultPrefix: "npm.js"
}
¡¡¡
`,
}
