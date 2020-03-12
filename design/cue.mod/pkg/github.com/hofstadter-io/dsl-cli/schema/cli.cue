package schema

Cli : {
  Name:     string
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
