package hof

import (
	"github.com/hofstadter-io/cuemod--cli-golang/schema"
)

AuthCommand :: schema.Command & {
  Name:    "auth"
  Usage:   "auth"
  Short:   "authentication subcommands"
  Long:    Short
  OmitRun: true

  Commands: [
    schema.Command & {
      Name:  "login"
      Usage: "login"
      Short: "login to an account"
      Long:  Short

      Imports: [
        {Path: "fmt"}
        // {Path: "github.com/hofstadter-io/mvs/lib"},
      ]

      Body: """
        fmt.Println("\(Name) login not implemented")
      """
    },
  ]
},

