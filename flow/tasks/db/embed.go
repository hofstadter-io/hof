package db

import (
  _ "embed"
  "fmt"

  "cuelang.org/go/cue"
)

//go:embed schema.cue
var task_schema string

var task_call cue.Value

func init_schemas (ctx *cue.Context) {
  if task_call.Exists() {
    return
  }

  val := ctx.CompileString(task_schema, cue.Filename("@embed:flow/tasks/db/schema.cue"))
  if val.Err() != nil {
    fmt.Println(val.Err())
    panic("should not have a schema error")
  }

  task_call = val.LookupPath(cue.ParsePath("Call"))
  if !task_call.Exists() {
    panic("missing flow/tasks/db.Call schema")
  }
}
