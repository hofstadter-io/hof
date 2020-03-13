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
          CMD: C & {
            PackageName: "commands"
          }
        }
      },
    }
    for _, C in Cli.Commands
  ]

  _SubCmds:  [[C & { Parent: P.In.CMD } for _, C in P.In.CMD.Commands] for _, P in _Commands]

  _SubCommands: [ // List comprehension
    {
      gen.CommandGen & {
        In: {
          CLI: Cli
          CMD: C
        }
      },
    }
    for _, C in list.FlattenN( _SubCmds, 1)
  ]

  _All: [_OnceFiles, _Commands, _SubCommands]
  Out: list.FlattenN(_All , 1)
}
