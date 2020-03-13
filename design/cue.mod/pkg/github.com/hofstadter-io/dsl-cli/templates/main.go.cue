package templates

MainTemplate : RealMainTemplate

RealMainTemplate : """
package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"{{ .CLI.Package }}/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
"""
