package runtime

func (ts *Script) CmdHello(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? stdin")
	}
	out := "hello world"
	if len(args) >= 1 {
		out = "hello " + args[0]
	}
	ts.stdout = out
}
