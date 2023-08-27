package middleware

import (
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	hofcontext "github.com/hofstadter-io/hof/flow/context"

	"github.com/hofstadter-io/hof/flow/middleware/info"
	"github.com/hofstadter-io/hof/flow/middleware/sync"
)

func UseDefaults(ctx *hofcontext.Context, opts flags.RootPflagpole, popts flags.FlowPflagpole) {
	// ctx.Use(dummy.NewDummy(opts, popts))
	ctx.Use(info.NewPrint(opts, popts))
	ctx.Use(info.NewProgress(opts, popts))
	//ctx.Use(info.NewBookkeeping(info.BookkeepingConfig{
	//Workdir: ".hof/flow",
	//}, opts, popts))
	ctx.Use(sync.NewPool(opts, popts))
}
