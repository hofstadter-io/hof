package commands

import (
	"fmt"

	// custom imports

	// infered imports
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/extern"
)

// Tool:   hof
// Name:   new
// Usage:  new <what> <name> <template>
// Parent: hof

var NewLong = `Make new apps, types, modules, pages, and functions.

<what>     - one of [module, type, page, component, func]
<name>     - the name for the new <what>
<template> - https://github.com/hofstadter-io/studios-new-templates@beta#modules/account-default

'<git-url>'    - a full url to a git repository
'@version'     - a branch, commit, or tag
'#nested-path' - a filepath within the repository
`

var (
	NewDataFlag string
)

func init() {
	NewCmd.Flags().StringVarP(&NewDataFlag, "data", "d", "", "a filepath or raw data in [json,xml,yaml,toml]")
	viper.BindPFlag("data", NewCmd.Flags().Lookup("data"))

}

var NewCmd = &cobra.Command{

	Use: "new <what> <name> <template>",

	Short: "Make new apps, types, modules, pages, and components.",

	Long: NewLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In newCmd", "args", args)
		// Argument Parsing
		// [0]name:   what
		//     help:
		//     req'd:  true
		if 0 >= len(args) {
			fmt.Println("missing required argument: 'what'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var what string

		if 0 < len(args) {

			what = args[0]
		}

		// [1]name:   name
		//     help:
		//     req'd:  true
		if 1 >= len(args) {
			fmt.Println("missing required argument: 'name'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var name string

		if 1 < len(args) {

			name = args[1]
		}

		// [2]name:   template
		//     help:
		//     req'd:

		var template string

		if 2 < len(args) {

			template = args[2]
		}

		/*
			fmt.Println("hof new:",
				what,

				name,

				template,
			)
		*/

		out, err := extern.NewEntry(what, name, template, NewDataFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(out)

	},
}

func init() {
	RootCmd.AddCommand(NewCmd)
}

func init() {
	// add sub-commands to this command when present

}
