package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/mod"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var modLong = `hof mod is CUE dependency management based on Go mods.

### Module File

The module file holds the requirements for project.
It is found in cue.mod/module.cue	

---
// These are like golang import paths
//   i.e. github.com/hofstadter-io/hof
module: "<module-path>"
cue: "v0.5.0"

// Required dependencies section
require: {
  // "<module-path>": "<module-semver>"
  "github.com/hofstadter-io/ghacue": "v0.2.0"
  "github.com/hofstadter-io/hofmod-cli": "v0.8.1"
}

// Indirect dependencies (managed by hof)
indirect: { ... }

// Replace dependencies with local or remote
replace: {
  "github.com/hofstadter-io/ghacue": "github.com/myorg/ghacue": "v0.4.2"
  "github.com/hofstadter-io/hofmod-cli": "../mods/clie"
}
---


### Authentication and private modules

hof mod prefers authenticated requests when fetching dependencies.
This increase rate limits with hosts and supports private modules.
Both token and sshkey base methods are supported, with preferences:

1. Matching entry in .netrc

2. ENV VARS for well known hosts.

  GITHUB_TOKEN
  GITLAB_TOKEN
  BITBUCKET_USERNAME / BITBUCKET_PASSWORD or BITBUCKET_TOKEN 

  The bitbucket method will depend on the account type and enterprise license.

3. SSH keys 

  the following are searched: ~/.ssh/config, /etc/ssh/config, ~/.ssh/in_rsa

  You can configure the SSH key with HOF_SSHUSR and HOF_SSHKEY


### Usage

there are two main commands you will use, init & tidy

# Initialize the current folder as a module
hof mod init <module-path>     (like github.com/org/repo)

# Refresh dependencies, discovering any new imports
hof mod tidy

# Add a dependency
hof mod get github.com/hofstadter-io/hof@v0.6.8
hof mod get github.com/hofstadter-io/hof@v0.6.8-beta.6
hof mod get github.com/hofstadter-io/hof@latest   // latest semver
hof mod get github.com/hofstadter-io/hof@next     // next prerelease
hof mod get github.com/hofstadter-io/hof@main     // latest commit on branch

# Update dependencies
hof mod get github.com/hofstadter-io/hof@latest
hof mod get all@latest

# Symlink dependencies from local cache
hof mod link

# Copy dependency code from local cache
hof mod vendor

# Verify dependency code against cue.mod/sums.cue
hof mod verify

# This helpful output
hof mod help

`

var ModCmd = &cobra.Command{

	Use: "mod",

	Aliases: []string{
		"m",
	},

	Short: "CUE dependency management based on Go mods",

	Long: modLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

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

	ModCmd.AddCommand(cmdmod.InitCmd)
	ModCmd.AddCommand(cmdmod.GetCmd)
	ModCmd.AddCommand(cmdmod.VerifyCmd)
	ModCmd.AddCommand(cmdmod.TidyCmd)
	ModCmd.AddCommand(cmdmod.LinkCmd)
	ModCmd.AddCommand(cmdmod.VendorCmd)

}
