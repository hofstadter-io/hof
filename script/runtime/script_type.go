package runtime

import (
	"bytes"
	"context"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/parnurzeal/gorequest"
	"go.uber.org/zap"

	"github.com/hofstadter-io/hof/lib/gotils/txtar"
)

// A Script holds execution state for a single test script.
type Script struct {
	// script params
	// Cfg           Config
	params Params

	ctxt context.Context // per Script context

	t           T
	testTempDir string
	workdir     string // temporary work dir ($WORK)
	cd          string // current directory during test execution; initially $WORK/gopath/src

	log bytes.Buffer // test execution log (printed at end of test)

	Logger *zap.SugaredLogger

	mark int // offset of next log truncation

	name   string // short name of test ("foo")
	file   string // full file name ("testdata/script/foo.hls")
	orig   string // original content
	lineno int    // line number currently executing
	line   string // line currently executing

	env    []string                    // environment list (for os/exec)
	envMap map[string]string           // environment mapping (matches env; on Windows keys are lowercase)
	values map[any]any // values for custom commands

	stdin  string // standard input to next 'go' command; set by 'stdin' command.
	stdout string // standard output from last 'go' command; for 'stdout' command
	stderr string // standard error from last 'go' command; for 'stderr' command
	status int    // status code from exec or http

	stopped bool      // test wants to stop early
	start   time.Time // time phase started
	failed  bool

	background []backgroundCmd // backgrounded 'exec' and 'go' commands
	deferred   func()          // deferred cleanup actions.

	archive       *txtar.Archive    // the testscript being run.
	scriptFiles   map[string]string // files stored in the txtar archive (absolute paths -> path in script)
	scriptUpdates map[string]string // updates to testscript files via UpdateScripts.

	fs billy.Filesystem

	fsClients   map[string]billy.Filesystem
	httpClients map[string]*gorequest.SuperAgent
	// cueRuntimes map[string]cue.Runtime

}

type Phase struct {
	Lines []string
	Steps []string
}
