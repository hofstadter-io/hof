package cmds

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
    {
      Name:  "login"
      Usage: "login"
      Short: "login to an account"
      Long:  Short
    },
    {
      Name:  "test"
      Usage: "test [name]"
      Short: "test your auth configuration, defaults to current"
      Long:  Short
    },
  ]
},

#ConfigCommand: schema.#Command & {
  Name:  "config"
  Usage: "config"
  Short: "configuration subcommands"
  Long:  Short

  OmitRun: true

  Commands: [
    {
      Name:  "create"
      Usage: "create"
      Short: "create a configuration"
      Long:  Short
    },
    {
      Name:  "list"
      Usage: "list"
      Short: "list configurations"
      Long:  Short
    },
    {
      Name:  "get"
      Usage: "get"
      Short: "print a configuration"
      Long:  Short
    },
    {
      Name:  "set"
      Usage:  "set"
      Short: "set configuration value(s)"
      Long:  Short
    },
    {
      Name:  "use"
      Usage: "use"
      Short: "set the default configuration"
      Long:  Short
    },
  ]
}

#SecretCommand: schema.#Command & {
  Name:  "secret"
  Usage: "secret"
  Short: "secret subcommands"
  Long:  Short

  OmitRun: true

  Commands: [
    {
      Name:  "create"
      Usage: "create"
      Short: "create secrets"
      Long:  Short
    },
    {
      Name:  "list"
      Usage: "list"
      Short: "list secrets"
      Long:  Short
    },
    {
      Name:  "get"
      Usage: "get"
      Short: "print a secret"
      Long:  Short
    },
    {
      Name:  "set"
      Usage:  "set"
      Short: "set secret value(s)"
      Long:  Short
    },
    {
      Name:  "use"
      Usage: "use"
      Short: "set the default configuration"
      Long:  Short
    },
  ]
}
