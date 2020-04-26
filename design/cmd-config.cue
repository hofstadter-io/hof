package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#ConfigCommand: schema.#Command & {
  Name:  "config"
  Usage: "config"
  Short: "configuration subcommands"
  Long:  Short

  OmitRun: true

  Commands: [
    schema.#Command & {
      Name:  "test"
      Usage: "test [name]"
      Short: "test your auth configuration, defaults to current"
      Long:  Short
      Args: [
        schema.#Arg & {
          Name:     "name"
          Type:     "string"
          Help:     "configuration name"
        },
      ]
    },
    schema.#Command & {
      Name:  "list"
      Usage: "list"
      Short: "list configurations"
      Long:  Short
    },
    schema.#Command & {
      Name:  "get"
      Usage: "get"
      Short: "print a configuration"
      Long:  Short
      Args: [
        schema.#Arg & {
          Name:     "name"
          Type:     "string"
          Required: true
          Help:     "name of the configuration"
        },
      ]
    },
    schema.#Command & {
      Name:  "set"
      Usage: "set <name> <host> <account> [project]"
      Short: "set configuration values"
      Long:  Short
      Args: [
        schema.#Arg & {
          Name:     "name"
          Type:     "string"
          Required: true
          Help:     "name for the configuration"
        },
        schema.#Arg & {
          Name:     "host"
          Type:     "string"
          Required: true
          Help:     "host for this configuration"
        },
        schema.#Arg & {
          Name:     "account"
          Type:     "string"
          Required: true
          Help:     "account for this configuration"
        },
        schema.#Arg & {
          Name: "project"
          Type: "string"
          Help: "default project for this configuration"
        },
      ]
    },
    schema.#Command & {
      Name:  "use"
      Usage: "use"
      Short: "set the default configuration"
      Long:  Short
      Args: [
        schema.#Arg & {
          Name:     "name"
          Type:     "string"
          Required: true
          Help:     "name of the configuration"
        },
      ]
    },
  ]
},
