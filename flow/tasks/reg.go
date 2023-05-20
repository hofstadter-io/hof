package tasks

import (
	hofcontext "github.com/hofstadter-io/hof/flow/context"

	"github.com/hofstadter-io/hof/flow/tasks/api"
	"github.com/hofstadter-io/hof/flow/tasks/csp"
	"github.com/hofstadter-io/hof/flow/tasks/cue"
	"github.com/hofstadter-io/hof/flow/tasks/db"
	"github.com/hofstadter-io/hof/flow/tasks/gen"
	"github.com/hofstadter-io/hof/flow/tasks/hof"
	"github.com/hofstadter-io/hof/flow/tasks/kv"
	"github.com/hofstadter-io/hof/flow/tasks/msg"
	"github.com/hofstadter-io/hof/flow/tasks/os"
	"github.com/hofstadter-io/hof/flow/tasks/st"
)

func RegisterDefaults(context *hofcontext.Context) {
	context.Register("noop", NewNoop)
	context.Register("nest", NewNest)

	context.Register("api.Call", api.NewCall)
	context.Register("api.Serve", api.NewServe)

	context.Register("csp.Chan", csp.NewChan)
	context.Register("csp.Recv", csp.NewRecv)
	context.Register("csp.Send", csp.NewSend)

	context.Register("cue.Format", cue.NewCueFormat)

	context.Register("db.Call", db.NewCall)

	context.Register("gen.CUID", gen.NewCUID)
	context.Register("gen.Float", gen.NewFloat)
	context.Register("gen.Int", gen.NewInt)
	context.Register("gen.Norm", gen.NewNorm)
	context.Register("gen.Now", gen.NewNow)
	context.Register("gen.Seed", gen.NewSeed)
	context.Register("gen.Slug", gen.NewSlug)
	context.Register("gen.Str", gen.NewStr)
	context.Register("gen.UUID", gen.NewUUID)

	context.Register("hof.Template", hof.NewHofTemplate)

	context.Register("kv.Mem", kv.NewMem)

	context.Register("msg.IrcClient", msg.NewIrcClient)

	context.Register("os.Exec", os.NewExec)
	context.Register("os.FileLock", os.NewFileLock)
	context.Register("os.FileUnlock", os.NewFileUnlock)
	context.Register("os.Getenv", os.NewGetenv)
	context.Register("os.Glob", os.NewGlob)
	context.Register("os.Mkdir", os.NewMkdir)
	context.Register("os.ReadFile", os.NewReadFile)
	context.Register("os.Sleep", os.NewSleep)
	context.Register("os.Stdin", os.NewStdin)
	context.Register("os.Stdout", os.NewStdout)
	context.Register("os.Watch", os.NewWatch)
	context.Register("os.WriteFile", os.NewWriteFile)

	context.Register("st.Diff", st.NewDiff)
	context.Register("st.Insert", st.NewInsert)
	context.Register("st.Mask", st.NewMask)
	context.Register("st.Patch", st.NewPatch)
	context.Register("st.Pick", st.NewPick)
	context.Register("st.Replace", st.NewReplace)
	context.Register("st.Upsert", st.NewUpsert)

}
