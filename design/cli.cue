package cli

import (
  "github.com/hofstadter-io/dsl-cli:cli"
  "github.com/hofstadter-io/dsl-cli/schema"
)

GEN : cli.Generator & {
  Cli: CLI
}

CLI : cli.Schema & {
  Name: "hof"
  Package: "github.com/hofstadter-io/hof"
  Short: "hof is a CLI"
  OmitRun: true

  Commands: [
    schema.Command & {
      Name: "auth"
    },
    schema.Command & {
      Name: "config"
      Commands: [
        schema.Command & {
          Name: "get"
        },
        schema.Command & {
          Name: "set"
        },
      ]
    },
  ]
}

