package os

import (
	_ "embed"
	"fmt"

	"cuelang.org/go/cue"
)

//go:embed schema.cue
var task_schema string

var task_exec cue.Value
// var task_chan cue.Value
// var task_send cue.Value
// var task_recv cue.Value
var task_watch cue.Value

func init_schemas(ctx *cue.Context) {
	if task_exec.Exists() {
		return
	}

	val := ctx.CompileString(task_schema, cue.Filename("@embed:flow/tasks/csp/schema.cue"))
	if val.Err() != nil {
		fmt.Println(val.Err())
		panic("should not have a schema error")
	}

	task_exec = val.LookupPath(cue.ParsePath("Exec"))
	if !task_exec.Exists() {
		panic("missing flow/tasks/os.Exec schema")
	}

	// task_chan = val.LookupPath(cue.ParsePath("Chan"))
	// if !task_chan.Exists() {
	// 	panic("missing flow/tasks/csp.Chan schema")
	// }
	// task_send = val.LookupPath(cue.ParsePath("Send"))
	// if !task_send.Exists() {
	// 	panic("missing flow/tasks/csp.Send schema")
	// }
	// task_recv = val.LookupPath(cue.ParsePath("Recv"))
	// if !task_recv.Exists() {
	// 	panic("missing flow/tasks/csp.Recv schema")
	// }

	task_watch = val.LookupPath(cue.ParsePath("Watch"))
	if !task_watch.Exists() {
		panic("missing flow/tasks/fs.Watch schema")
	}
}
