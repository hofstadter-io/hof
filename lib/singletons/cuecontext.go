package singletons

import (
	"sync"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

var cueContext *cue.Context
var cueContextMutex sync.Mutex

func init() {
	cueContext = cuecontext.New()
}

func CueContext() *cue.Context {
	return cueContext
}

func EmptyValue() cue.Value {
	cueContextMutex.Lock()
	defer cueContextMutex.Unlock()

	return cueContext.CompileString("{}")
}

func CompileString(src string, opts ...cue.BuildOption) cue.Value {
	cueContextMutex.Lock()
	defer cueContextMutex.Unlock()

	return cueContext.CompileString(src, opts...)
}

func CompileBytes(src []byte, opts ...cue.BuildOption) cue.Value {
	cueContextMutex.Lock()
	defer cueContextMutex.Unlock()

	return cueContext.CompileBytes(src, opts...)
}
