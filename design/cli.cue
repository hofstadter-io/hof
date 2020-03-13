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

  Commands: [
    schema.Command & {
      Name: "auth"
    },
    schema.Command & {
      Name: "config"
    },
  ]
}

