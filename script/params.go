package script

// Params holds parameters for a call to Run.
type Params struct {
	// Mode is the top level operation mode for script
	// It can be "test" or "run" to control how scripts are isolated (or not)
	Mode string

	// Dir holds the name of the directory holding the scripts.
	// All files in the directory with a .hls suffix will be considered
	// as test scripts. By default the current directory is used.
	// Dir is interpreted relative to the current test directory.
	Dir string

	// Glob holds the patter to match, defaults to '*.hls'
	Glob string

	// Setup is called, if not nil, to complete any setup required
	// for a test. The WorkDir and Vars fields will have already
	// been initialized and all the files extracted into WorkDir,
	// and Cd will be the same as WorkDir.
	// The Setup function may modify Vars and Cd as it wishes.
	Setup func(*Env) error

	// Condition is called, if not nil, to determine whether a particular
	// condition is true. It's called only for conditions not in the
	// standard set, and may be nil.
	Condition func(cond string) (bool, error)

	// Cmds holds a map of commands available to the script.
	// It will only be consulted for commands not part of the standard set.
	Cmds map[string]func(ts *Script, neg int, args []string)

	// Funcs holds a map of functions available to the script.
	// These work like exec and use 'call' instead.
	// Use these to facilitate code coverage (exec does not capture this).
	Funcs map[string]func(ts *Script, args []string) error

	// TestWork specifies that working directories should be
	// left intact for later inspection.
	TestWork bool

	// WorkdirRoot specifies the directory within which scripts' work
	// directories will be created. Setting WorkdirRoot implies TestWork=true.
	// If empty, the work directories will be created inside
	// $GOTMPDIR/go-test-script*, where $GOTMPDIR defaults to os.TempDir().
	WorkdirRoot string

	// IgnoreMissedCoverage specifies that if coverage information
	// is being generated (with the -test.coverprofile flag) and a subcommand
	// function passed to RunMain fails to generate coverage information
	// (for example because the function invoked os.Exit), then the
	// error will be ignored.
	IgnoreMissedCoverage bool

	// UpdateScripts specifies that if a `cmp` command fails and
	// its first argument is `stdout` or `stderr` and its second argument
	// refers to a file inside the testscript file, the command will succeed
	// and the testscript file will be updated to reflect the actual output.
	//
	// The content will be quoted with txtar.Quote if needed;
	// a manual change will be needed if it is not unquoted in the
	// script.
	UpdateScripts bool

	// Line prefix which indicates a new phase
	// defaults to "#"
	PhasePrefix string

	// Comment prefix for a line
	// defaults to "~"
	CommentPrefix string
}

func paramDefaults(p Params) Params {
	if p.Mode == "" {
		p.Mode = "test"
	}
	if p.Glob == "" {
		p.Glob = "*.hls"
	}
	if p.PhasePrefix == "" {
		p.PhasePrefix = "###"
	}
	if p.CommentPrefix == "" {
		p.CommentPrefix = "#"
	}

	return p
}

