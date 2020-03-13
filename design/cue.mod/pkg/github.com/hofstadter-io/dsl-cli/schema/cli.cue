package schema

import (
  "strings"
)

Cli : {
  Name:     string
  cliName: strings.ToCamel(Name)
  CliName: strings.ToTitle(Name)

  Package:  string

  Usage?:    string
  Short?:    string
  Long?:     string

  PersistentPrerun?:   bool | *false
  Prerun?:             bool | *false
  OmitRun?:            bool | *false
  Postrun?:            bool | *false
  PersistentPostrun?:  bool | *false

  Pflags?:    [...Flag]
  Flags?:     [...Flag]
  Args?:      [...Arg]
  Commands?:  [...Command]
}
