package runtime

// regexp checks that file content matches a regexp.
// it accepts Go regexp syntax.
func (ts *Script) CmdRegexp(neg int, args []string) {
	scriptMatch(ts, neg, args, "", "regexp")
}

// regexp checks that file content matches a regexp.
// it accepts Go regexp syntax and returns the matches
func (ts *Script) CmdGrep(neg int, args []string) {
	scriptMatch(ts, neg, args, "", "grep")
}

// sed finds and replaces in text content
// it accepts Go regexp syntax and returns the replaced content
func (ts *Script) CmdSed(neg int, args []string) {
	scriptMatch(ts, neg, args, "", "sed")
}
