package function

import (
	"fmt"

	// custom imports

	// infered imports
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/fns"
)

// Tool:   hof
// Name:   create
// Usage:  create [path/to]<name> <template>[@version][#template-subpath]
// Parent: function

var CreateLong = `Create a new function from a template. The path prefix says where, the last part will be the name`

var (
	CreateHereFlag bool

	CreateTemplateFlag string
)

func init() {
	CreateCmd.Flags().BoolVarP(&CreateHereFlag, "here", "h", false, "create in the current directory (uses dir as name)")
	viper.BindPFlag("here", CreateCmd.Flags().Lookup("here"))

	CreateCmd.Flags().StringVarP(&CreateTemplateFlag, "template", "t", "https://github.com/hofstadter-io/studios-functions#custom-default", "create with a template, set to empty '-t' to omit dir/file creation")
	viper.BindPFlag("template", CreateCmd.Flags().Lookup("template"))

}

var CreateCmd = &cobra.Command{

	Use: "create [path/to]<name> <template>[@version][#template-subpath]",

	Short: "Create a new function",

	Long: CreateLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In createCmd", "args", args)
		// Argument Parsing
		// [0]name:   name
		//     help:
		//     req'd:  true
		if 0 >= len(args) {
			fmt.Println("missing required argument: 'name'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var name string

		if 0 < len(args) {

			name = args[0]
		}

		/*
			fmt.Println("hof function create:",
				name,

				template,
			)
		*/

		err := fns.Create(name, CreateTemplateFlag, CreateHereFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	// add sub-commands to this command when present

}
