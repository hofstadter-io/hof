package runtime

func (ts *Script) CmdStdin(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? stdin")
	}
	if len(args) != 1 {
		ts.Fatalf("usage: stdin filename")
	}
	ts.stdin = ts.ReadFile(args[0])
}

// stdout checks that the last go command standard output matches a regexp.
func (ts *Script) CmdStdout(neg int, args []string) {
	scriptMatch(ts, neg, args, ts.stdout, "stdout")
}

// stderr checks that the last go command standard output matches a regexp.
func (ts *Script) CmdStderr(neg int, args []string) {
	scriptMatch(ts, neg, args, ts.stderr, "stderr")
}
