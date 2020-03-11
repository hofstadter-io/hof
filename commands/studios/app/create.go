package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hofstadter-io/hof/pkg/studios/app"
)

var CreateLong = `Create an app from a template or existing directory`

var (
	CreateHereFlag bool

	CreateTemplateFlag string
)

func init() {
	CreateCmd.Flags().BoolVarP(&CreateHereFlag, "here", "", false, "create in the current directory (uses dir as name)")
	viper.BindPFlag("here", CreateCmd.Flags().Lookup("here"))

	CreateCmd.Flags().StringVarP(&CreateTemplateFlag, "template", "t", "https://github.com/hofstadter-io/hof-starter-app", "create with a template, specifiying the 'here' flag will omit intial code creation")
	viper.BindPFlag("template", CreateCmd.Flags().Lookup("template"))

}

var CreateCmd = &cobra.Command{

	Use: "create <name> <app-version> <template>[@version]",

	Short: "Create an app.",

	Long: CreateLong,

	Run: func(cmd *cobra.Command, args []string) {

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
			fmt.Println("hof app create:",
				name,
			)
		*/

		err := app.Create(name, CreateTemplateFlag, CreateHereFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
