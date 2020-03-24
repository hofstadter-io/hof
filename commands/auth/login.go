package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

var loginLong = `login to an account`

var LoginCmd = &cobra.Command{

	Use: "login",

	Short: "login to an account",

	Long: loginLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		fmt.Println("hof login login not implemented")

	},
}
