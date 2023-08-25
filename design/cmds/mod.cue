package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

ModCmdImports: [
	{Path: "github.com/hofstadter-io/hof/lib/mod", ...},
	{Path: "github.com/hofstadter-io/hof/cmd/hof/flags", ...},
]

ModCommand: schema.Command & {
	// TBD:   "β"
	Name:  "mod"
	Usage: "mod"
	Aliases: ["m"]
	Short: "CUE module dependency management"
	Long:  ModLongHelp

	//Topics: #ModTopics
	//Examples: #ModExamples

	OmitRun: true

	#body: {
		func:    string
		module:  bool | *false
		_modstr: string | *""
		if module == true {
			_modstr: "module, "
		}
		content: """
			err = mod.\(func)(\(_modstr) flags.RootPflags)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			"""
	}

	Commands: [{
		Name:  "init"
		Usage: "init <module>"
		Short: "initialize a new module in the current directory"
		Long:  Short

		Args: [{
			Name:     "module"
			Type:     "string"
			Required: true
			Help:     "module path"
		}]

		Imports: ModCmdImports
		Body:    (#body & {func: "Init", module: true}).content
	}, {
		Name:  "get"
		Usage: "get <module>"
		Short: "add a new dependency to the current module"
		Long:  Short

		Args: [{
			Name:     "module"
			Type:     "string"
			Required: true
			Help:     "module path@version"
		}]

		Flags: [...schema.Flag] & [ {
			Name:    "Prerelease"
			Long:    "prerelease"
			Short:   "P"
			Type:    "bool"
			Default: "false"
			Help:    "include prerelease version when using @latest"
		}]

		Imports: [
			{Path: "github.com/hofstadter-io/hof/lib/mod", ...},
		]

		Body: """
			err = mod.Get(module, flags.RootPflags, flags.Mod__GetFlags)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			"""
	}, {
		Name:  "verify"
		Usage: "verify"
		Short: "verify integrity of dependencies"
		Long:  Short

		Imports: ModCmdImports

		Body: (#body & {func: "Verify"}).content
	}, {
		Name:  "tidy"
		Usage: "tidy"
		Short: "recalculate dependencies and update mod files"
		Long:  Short

		Imports: ModCmdImports

		Body: (#body & {func: "Tidy"}).content
	}, {
		Name:  "link"
		Usage: "link"
		Short: "symlink dependencies to cue.mod/pkg"
		Long:  Short

		Imports: ModCmdImports

		Body: (#body & {func: "Link"}).content
	}, {
		Name:  "vendor"
		Usage: "vendor"
		Short: "copy dependencies to cue.mod/pkg"
		Long:  Short

		Imports: ModCmdImports

		Body: (#body & {func: "Vendor"}).content
	}, {
		Name:  "clean"
		Usage: "clean"
		Short: "clean hof's module cache"
		Long:  Short

		Imports: ModCmdImports

		Body: (#body & {func: "Clean"}).content
	}, {
		Name:  "publish"
		Usage: "publish <module>"
		Short: "publish a module"
		Long:  Short

		Args: [{
			Name:     "module"
			Type:     "string"
			Required: true
			Help:     "module path@version"
		}]

		Imports: [
			{Path: "github.com/hofstadter-io/hof/lib/mod", ...},
		]

		Body: """
			if err = mod.Publish(module); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			"""
	}]

}

ModLongHelp: string & ##"""
	hof mod is CUE dependency management based on Go mods.
	
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
	  machine github.com
	  login github-token
	  password <github-token-value>
	
	2. ENV VARS for well known hosts.
	
	  GITHUB_TOKEN
	  GITLAB_TOKEN
	  BITBUCKET_USERNAME / BITBUCKET_PASSWORD or BITBUCKET_TOKEN 

	  The bitbucket method will depend on the account type and enterprise license.
	
	3. SSH keys 

	  the following are searched: ~/.ssh/config, /etc/ssh/config, ~/.ssh/id_rsa
	
	
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
	
	
	"""##
