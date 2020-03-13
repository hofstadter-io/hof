package schema

import (
  "strings"
)

Command : {
  Name:     string
  cmdName: strings.ToCamel(Name)
  CmdName: strings.ToTitle(Name)

  Aliases?: [...string]
  PackageName?: string

  // TODO, generate usage, and maybe long help
  Usage?:    string
  Short?:    string
  Long?:     string

  PersistentPrerun?:   bool | *false
  Prerun?:             bool | *false
  OmitRun?:            bool | *false
  Postrun?:            bool | *false
  PersistentPostrun?:  bool | *false

  Pflags?:   [...Flag]
  Flags?:    [...Flag]
  Args?:     [...Arg]
  Commands:  [...Command] | *[]
}

