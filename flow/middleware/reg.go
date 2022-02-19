package middleware

import (
  hofcontext "github.com/hofstadter-io/hof/flow/context"
  
	"github.com/hofstadter-io/hof/flow/middleware/dummy"
)

func ApplyDefaults(context *hofcontext.Context) {
  context.Apply(&dummy.Dummy{})
}
