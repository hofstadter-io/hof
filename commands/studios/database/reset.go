package db

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hofstadter-io/hof/pkg/studios/db"
)

var ResetLong = `Resets the DB to the latest checkpoint & schema, also adding seed data.`

var (
	ResetHardFlag bool
)

func init() {
	ResetCmd.Flags().BoolVarP(&ResetHardFlag, "hard", "H", false, "perform a hard database rest, squashing all migrations into one.")
	viper.BindPFlag("hard", ResetCmd.Flags().Lookup("hard"))

}

var ResetCmd = &cobra.Command{

	Use: "reset",

	Short: "Reset the DB to the latest schema",

	Long: ResetLong,

	Run: func(cmd *cobra.Command, args []string) {

		// fmt.Println("hof db reset:")

		err := db.Reset(ResetHardFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
