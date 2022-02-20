package middleware

import ( 
	"github.com/hofstadter-io/hof/cmd/hof/flags"
  hofcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/middleware/dummy"
	"github.com/hofstadter-io/hof/flow/middleware/info"
)

func UseDefaults(ctx *hofcontext.Context, opts *flags.RootPflagpole, popts *flags.FlowFlagpole) {
  ctx.Use(dummy.NewDummy(opts, popts))
  ctx.Use(info.NewProgress(opts, popts))
  ctx.Use(info.NewBookkeeping(opts, popts))
}
