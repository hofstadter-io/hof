package runtime

import (
	"context"
	"io"

	"github.com/go-git/go-billy/v5"
	"github.com/parnurzeal/gorequest"
	"go.uber.org/zap"

	"github.com/hofstadter-io/hof/script/ast"
)

type Runtime struct {
	// normal stuff
	params *Params
	result *ast.Result

	// testing related
	t      T

	// bookkeeping helpers
	parser  *ast.Parser
	script  *ast.Script
	phase   *ast.Phase
	currcmd *ast.Cmd
	lastcmd *ast.Cmd

	// background    []backgroundCmd             // backgrounded 'exec' and 'go' commands
	deferred      func()                      // deferred cleanup actions.

	stdoutStr string
	stderrStr string
	status  int
	failed  bool                        // runtime wants to stop early with failed state
	stopped bool                        // runtime wants to stop early

	// input, output, loggin, context helpers
	ctxt   context.Context
	logger *zap.SugaredLogger
	stdinR io.ReadCloser
	stdinW io.Writer
	stdout io.Writer
	stderr io.Writer
	// multi-writers, typically
	StdinW io.Writer
	Stdout io.Writer
	Stderr io.Writer

	// dir helpers
	currdir string
	testdir string
	workdir string

	// maps
	envMap      map[string]string
	values      map[string]interface{}
	aliases     map[string]string
	fsClients   map[string]billy.Filesystem
	httpClients map[string]*gorequest.SuperAgent
	// cueRuntimes map[string]cue.Runtime
	// k8sRuntimes map[string]cue.Runtime
	// file descriptors?

	// need to understand context for commands like gcloud (default)
}

func (RT *Runtime) SetMultiWriters(R *ast.Result) {
	// RT.Stdin = io.MultiWriter(RT.stdin, R.stdin)
	RT.Stdout = io.MultiWriter(RT.stdout, R.Stdout)
	RT.Stderr = io.MultiWriter(RT.stderr, R.Stderr)
}

