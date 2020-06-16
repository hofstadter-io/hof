package runtime

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// call runs the given function.
func (ts *Script) CmdCall(neg int, args []string) {
	if len(args) < 1 {
		ts.Fatalf("usage: call function [args...]")
	}

	var err error
	ts.stdout, ts.stderr, err = ts.call(args[0], args[1:]...)
	if ts.stdout != "" {
		fmt.Fprintf(&ts.log, "[stdout]\n%s", ts.stdout)
	}
	if ts.stderr != "" {
		fmt.Fprintf(&ts.log, "[stderr]\n%s", ts.stderr)
	}
	if err == nil && neg > 0 {
		ts.Fatalf("unexpected command success")
	}

	if err != nil {
		fmt.Fprintf(&ts.log, "[%v]\n", err)
		if ts.ctxt.Err() != nil {
			ts.Fatalf("test timed out while running command")
		} else if neg > 0 {
			ts.Fatalf("unexpected command failure")
		}
	}
}


// call runs the given function and then returns collected standard output and standard error.
func (ts *Script) call(function string, args ...string) (string, string, error) {

	fn, ok := ts.params.Funcs[function]
	if !ok {
		ts.Fatalf("unknown function %q", function)
		return "", "", fmt.Errorf("unknown function%q", function)
	}

	// backup originals
	oldstdout := os.Stdout
	oldstderr := os.Stderr
	stdout, outw, _ := os.Pipe()
	stderr, errw, _ := os.Pipe()
	os.Stdout = stdout
	os.Stderr = stderr

	var err error
	done := make(chan string)
	outC := make(chan string, 1)
	errC := make(chan string, 1)

	// call the function
	go func() {
		err = fn(ts, args)
		done <- "done"
	}()

	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var bufout bytes.Buffer
		io.Copy(&bufout, stdout)
		outC <- bufout.String()
	}()
	go func() {
		var buferr bytes.Buffer
		io.Copy(&buferr, stderr)
		errC <- buferr.String()
	}()

	// wait for function

	<-done
	// restore OS stds
	outw.Close()
	errw.Close()
	os.Stdout = oldstdout
	os.Stderr = oldstderr

	// get content
	funcout := <-outC
	funcerr := <-errC

	return funcout, funcerr, err
}

