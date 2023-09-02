package runtime

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/hofstadter-io/hof/lib/gotils/diff"
)

// cmp compares two files.
func (ts *Script) CmdCmp(neg int, args []string) {
	if neg != 0 {
		// It would be strange to say "this file can have any content except this precise byte sequence".
		ts.Fatalf("unsupported: !? cmp")
	}
	if len(args) < 2 {
		ts.Fatalf("usage: cmp file1 file2 [-trim-space]")
	}

	trim := false
	if len(args) > 2 {
		if args[2] == "-trim-space" {
			trim = true
		}
	}

	ts.doCmdCmp(args, false, trim)
}

// cmpenv compares two files with environment variable substitution.
func (ts *Script) CmdCmpenv(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? cmpenv")
	}
	if len(args) != 2 {
		ts.Fatalf("usage: cmpenv var1 var2")
	}
	ts.doCmdCmp(args, true, false)
}

func (ts *Script) doCmdCmp(args []string, env bool, trim bool) {
	// TODO, check args here

	arg2isStr := false

	// why are args (text1&2) handled differently here?
	// (later, answer is because you can check for existance of a string (arg2) in a file / output (arg1)
	text1 := ts.ReadFile(args[0])
	var text2 string
	// first try file, otherwise assume string
	absName2 := ts.MkAbs(args[1])
	_, err := os.Lstat(absName2)
	if err != nil {
		// file does not exist, assume it is a string
		if _, ok := err.(*os.PathError); ok {
			arg2isStr = true
			text2 = args[1]
		} else {
			ts.Check(err)
		}
	} else {
		data, err := os.ReadFile(absName2)
		ts.Check(err)
		text2 = string(data)
	}
	if env {
		text2 = ts.expand(text2)
	}

	// for strings, we just look to see if it exists anywhere, not complete match
	if (arg2isStr && strings.Contains(text1, text2)) || text1 == text2 {
		return
	}

	if ts.params.UpdateScripts && !env && (args[0] == "stdout" || args[0] == "stderr") {
		if scriptFile, ok := ts.scriptFiles[absName2]; ok {
			ts.scriptUpdates[scriptFile] = text1
			return
		}
		// The file being compared against isn't in the txtar archive, so don't
		// update the script.
	}

	// trim trailing whitespace
	if trim {
		var t1, t2 strings.Builder
		for _, line := range strings.Split(text1, "\n") {
			fmt.Fprintln(&t1, strings.TrimRightFunc(line, unicode.IsSpace))
		}
		for _, line := range strings.Split(text2, "\n") {
			fmt.Fprintln(&t2, strings.TrimRightFunc(line, unicode.IsSpace))
		}
		text1 = t1.String()
		text2 = t2.String()
	}

	// this is a separate trim to remove space from ends of file
	text1 = strings.TrimSpace(text1)
	text2 = strings.TrimSpace(text2)
	// fmt.Println(text1)

	if text1 == text2 {
		return
	}

	ts.Logf("[diff -%s(%d) +%s(%d)]\n%s\n", args[0], len(text1), args[1], len(text2), diff.Diff(args[0], []byte(text1), args[1], []byte(text2)))
	ts.Fatalf("%s and %s differ", args[0], args[1])
}
