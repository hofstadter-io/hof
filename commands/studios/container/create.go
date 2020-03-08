package container

import (
	"fmt"
	"os"

	// custom imports

	// infered imports

	"github.com/hofstadter-io/hof/lib/crun"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   create
// Usage:  create [path/to]<name> <template>[@version][#template-subpath]
// Parent: container

var CreateLong = `Create a new crun from a template. The path prefix says where, the last part will be the name`

var (
	CreateHereFlag bool

	CreateTemplateFlag string
)

func init() {
	CreateCmd.Flags().BoolVarP(&CreateHereFlag, "here", "", false, "create in the current directory (uses dir as name)")
	viper.BindPFlag("here", CreateCmd.Flags().Lookup("here"))

	CreateCmd.Flags().StringVarP(&CreateTemplateFlag, "template", "t", "https://github.com/hofstadter-io/studios-containers#custom-default", "create with a template, set to empty '-t' to omit dir/file creation")
	viper.BindPFlag("template", CreateCmd.Flags().Lookup("template"))

}

var CreateCmd = &cobra.Command{

	Use: "create [path/to]<name> <template>[@version][#template-subpath]",

	Short: "Create a new crun",

	Long: CreateLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In createCmd", "args", args)
		// Argument Parsing
		// [0]name:   name
		//     help:
		//     req'd:

		var name string

		if 0 < len(args) {

			name = args[0]
		}

		/*
			fmt.Println("hof containers create:",
				name,
			)
		*/

		err := crun.Create(name, CreateHereFlag, CreateTemplateFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// add sub-commands to this command when present

}
