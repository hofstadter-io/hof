package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#AuthCommand: schema.#Command & {
  Name:    "auth"
  Usage:   "auth"
  Short:   "authentication subcommands"
  Long:    Short
  OmitRun: true

  Commands: [
    schema.#Command & {
      Name:  "login"
      Usage: "login"
      Short: "login to an account"
      Long:  Short

    },
  ]
},

