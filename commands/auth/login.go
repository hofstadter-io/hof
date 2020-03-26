package auth

import (

	// hello... something might need to go here

	"github.com/spf13/cobra"

	"fmt"
)

var loginLong = `login to an account`

var LoginCmd = &cobra.Command{

	Use: "login",

	Short: "login to an account",

	Long: loginLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		fmt.Println("login login not implemented")

	},
}
