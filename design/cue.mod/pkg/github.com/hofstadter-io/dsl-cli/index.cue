package cli

import (
	"list"

  "github.com/hofstadter-io/dsl-cli/gen"
  "github.com/hofstadter-io/dsl-cli/schema"
)

Schema : schema.Cli

Generator : {
  _Cli: schema.Cli
  _OnceFiles: [
    gen.TestGen & {
      _In: {
        CLI: _Cli
      }
    },
	]
	_Commands: [ {
			gen.MultiGen & {
				_In: {
					CLI: _Cli
					CMD: C
				}
			},
		} for i, C in _Cli.Commands ]

	_All: [_OnceFiles, _Commands]
	_Out: _All

	// Having issues with flattening lists here
	// _Out: list.FlattenN(_All , 1)
}
