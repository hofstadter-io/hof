package cli

import (
	"list"

  "github.com/hofstadter-io/dsl-cli/gen"
  "github.com/hofstadter-io/dsl-cli/schema"
)

Schema : schema.Cli

Generator : {
  Cli: schema.Cli
  _OnceFiles: [
    gen.MainGen & {
      In: {
        CLI: Cli
      }
    },
    gen.RootGen & {
      In: {
        CLI: Cli
      }
    },
  ]
  _Commands: [ // List comprehension
    {
      gen.CommandGen & {
        In: {
          CLI: Cli
          CMD: C
        }
      },
    }
    for i, C in Cli.Commands
  ]

  _All: [_OnceFiles, _Commands]
  Out: list.FlattenN(_All , 1)
}
