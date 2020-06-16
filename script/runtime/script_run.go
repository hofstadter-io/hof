package runtime

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// run runs the test script.
func (ts *Script) run() {

	defer func() {
		ts.cleanup()
	}()

	var script string
	if ts.params.Mode == "test" {
		script = ts.setupTest()
	} else if ts.params.Mode == "run" {
		script = ts.setupRun()
	} else {
		fmt.Errorf("Unknown HLS mode in Params: %q", ts.params.Mode)
		return
	}

	// With -v or -testwork, start log with full environment.
	if *testWork || ts.t.Verbose() {
		if ts.params.Mode == "test" {
			// Display environment.
			ts.CmdEnv(0, nil)
			fmt.Fprintf(&ts.log, "\n")
			ts.mark = ts.log.Len()
		}
	}
	defer ts.applyScriptUpdates()

	// Run script.
	// See testdata/script/README for documentation of script form.
Script:
	for script != "" {
		// Extract next line.
		ts.lineno++
		var line string
		if i := strings.Index(script, "\n"); i >= 0 {
			line, script = script[:i], script[i+1:]
		} else {
			line, script = script, ""
		}

		// # is a comment indicating the start of new phase.
		if strings.HasPrefix(line, ts.params.PhasePrefix) {
			// If there was a previous phase, it succeeded,
			// so rewind the log to delete its details (unless -v is in use).
			// If nothing has happened at all since the mark,
			// rewinding is a no-op and adding elapsed time
			// for doing nothing is meaningless, so don't.
			if ts.log.Len() > ts.mark {
				ts.rewind()
				ts.markTime()
			}
			// Print phase heading and mark start of phase output.
			fmt.Fprintf(&ts.log, "%s\n", line)
			ts.mark = ts.log.Len()
			ts.start = time.Now()
			continue
		}

		// Line comment, this will be annoying to anyone already using testsuite
		// .. they can do a find and replace pretty easily though, can probably do with just sed
		if strings.HasPrefix(line, ts.params.CommentPrefix) {
			continue
		}

		// Parse input line. Ignore blanks entirely.
		args := ts.parse(line)
		if len(args) == 0 {
			continue
		}

		// Echo command to log.
		fmt.Fprintf(&ts.log, "\n>>> %s\n", line)

		// Command prefix [cond] means only run this command if cond is satisfied.
		for strings.HasPrefix(args[0], "[") && strings.HasSuffix(args[0], "]") {
			cond := args[0]
			cond = cond[1 : len(cond)-1]
			cond = strings.TrimSpace(cond)
			args = args[1:]
			if len(args) == 0 {
				ts.Fatalf("missing command after condition")
			}
			want := true
			if strings.HasPrefix(cond, "!") {
				want = false
				cond = strings.TrimSpace(cond[1:])
			}
			ok, err := ts.condition(cond)
			if err != nil {
				ts.Fatalf("bad condition %q: %v", cond, err)
			}
			if ok != want {
				// Don't run rest of line.
				continue Script
			}
		}

		// Command prefix ! means negate the expectations about this command:
		// go command should fail, match should not be found, etc.
		neg := 0
		if args[0] == "!" {
			neg = 1
			args = args[1:]
			if len(args) == 0 {
				ts.Fatalf("! on line by itself")
			}
		} else if args[0] == "?" {
			neg = -1
			args = args[1:]
			if len(args) == 0 {
				ts.Fatalf("? on line by itself")
			}
		}

		/* Run command, check order
			1. params commands, incase of overrides
			2. buildin commands
			3. fallback on exec ...
		*/

		// the command name
		C := args[0]

		// try user commands
		cmd, cmdOK := ts.params.Cmds[C]

		// check user command, try builtin command
		if !cmdOK || cmd == nil {
			cmd, cmdOK = scriptCmds[C]
		}

		// check builtin command, try system command
		if !cmdOK || cmd == nil {
			path, err := exec.LookPath(C)
			if err == nil {
				cmd = scriptCmds["exec"]
				nargs := []string{ C, path}
				args = append(nargs, args[1:]...)
			}
		}

		if cmd == nil {
			ts.Fatalf("unknown command or function %q", args[0])
		}
		cmd(ts, neg, args[1:])

		// Command can ask script to stop early.
		if ts.stopped {
			// Break instead of returning, so that we check the status of any
			// background processes and print PASS.
			break
		}
	}

	for _, bg := range ts.background {
		interruptProcess(bg.cmd.Process)
	}
	ts.CmdWait(0, nil)

	// Final phase ended.
	ts.rewind()
	ts.markTime()
	if !ts.stopped {
		if ts.params.Mode == "test" {
			fmt.Fprintf(&ts.log, "PASS\n")
		}
		if ts.params.Mode == "run" {
			fmt.Fprintf(&ts.log, "DONE\n")
		}
	}
}


// Truncate log at end of last phase marker,
// discarding details of successful phase.
func (ts *Script) rewind() {
	if !ts.t.Verbose() {
		ts.log.Truncate(ts.mark)
	}
}

func (ts *Script) markTime() {
	if ts.mark > 0 && !ts.start.IsZero() {
		afterMark := append([]byte{}, ts.log.Bytes()[ts.mark:]...)
		ts.log.Truncate(ts.mark - 1) // cut \n and afterMark
		fmt.Fprintf(&ts.log, " (%.3fs)\n", time.Since(ts.start).Seconds())
		ts.log.Write(afterMark)
	}
	ts.start = time.Time{}
}

// Insert elapsed time for phase at end of phase marker
func (ts *Script) cleanup() {

	// On a normal exit from the test loop, background processes are cleaned up
	// before we print PASS. If we return early (e.g., due to a test failure),
	// don't print anything about the processes that were still running.
	for _, bg := range ts.background {
		interruptProcess(bg.cmd.Process)
	}
	for _, bg := range ts.background {
		<-bg.wait
	}
	ts.background = nil

	ts.markTime()
	// Flush testScript log to testing.T log.
	ts.t.Log("\n" + ts.abbrev(ts.log.String()))
	ts.deferred()
}
