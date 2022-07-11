package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/mod"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/mod"
)

var modLong = `hof mod is a flexible tool and library based on Go mods.

Use and create module systems with Minimum Version Selection (MVS) semantics
and manage dependencies go mod style. Mix any set of language, code bases,
git repositories, package managers, and subdirectories.


### Features

- Based on go mods MVS system, aiming for 100% reproducible builds.
- Recursive dependencies, version resolution, and code instrospection.
- Custom module systems with custom file names and vendor directories.
- Control configuration for naming, vendoring, and other behaviors.
- Polyglot support for multiple module systems and multiple languages within one tool.
- Works with any git system and supports the main features from go mods.
- Convert other vendor and module systems to MVS or just manage their files with MVS.
- Private repository support for GitHub, GitLab, Bitbucket, and git+SSH.


### Usage

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


### Module File

The module file holds the requirements for project.
It has the same format as a go.mod file.

---
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
---


### Authentication and private modules

hof mod prefers authenticated requests when fetching dependencies.
This increase rate limits with hosts and supports private modules.
Both token and sshkey base methods are supported.

If you are using credentials, then private modules can be transparent.
An ENV VAR like GOPRIVATE and CUEPRIVATE can be used to require credentials.

The following ENV VARS are used to set credentials.

GITHUB_TOKEN
GITLAB_TOKEN
BITBUCKET_TOKEN or BITBUCKET_USERNAME / BITBUCKET_PASSWORD *

SSH keys will be looked up in the following ~/.ssh/config, /etc/ssh/config, ~/.ssh/in_rsa

You can configure the SSH key with

HOF_SSHUSR and HOF_SSHKEY

* The bitbucket method will depend on the account type and enterprise license.


### Custom Module Systems

.mvsconfig.cue allows you to define custom module systems.
With some simple configuration, you can create and control
and vendored module system based on go mods.
You can also define global configurations.

See the ./lib/mod/langs in the repository for examples.

### Motivation

- MVS has better semantics for vendoring and gets us closer to 100% reproducible builds.
- JS and Python can have MVS while still using the remainder of the tool chains.
- Prototype for cuelang module and vendor management.
- We need a module system for our [hof-lang](https://hof-lang.org) project.
`

func ModPersistentPreRun(args []string) (err error) {

	mod.InitLangs()

	return err
}

var ModCmd = &cobra.Command{

	Use: "mod",

	Aliases: []string{
		"m",
	},

	Short: "polyglot dependency management based on go mods and MVS",

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

	},
}

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := ModCmd.HelpFunc()
	ousage := ModCmd.UsageFunc()
	help := func(cmd *cobra.Command, args []string) {
		if extra(cmd) {
			return
		}
		ohelp(cmd, args)
	}
	usage := func(cmd *cobra.Command) error {
		if extra(cmd) {
			return nil
		}
		return ousage(cmd)
	}

	ModCmd.SetHelpFunc(help)
	ModCmd.SetUsageFunc(usage)

	ModCmd.AddCommand(cmdmod.InfoCmd)
	ModCmd.AddCommand(cmdmod.InitCmd)
	ModCmd.AddCommand(cmdmod.VendorCmd)

}