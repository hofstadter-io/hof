package runtime

import (
	"fmt"
	"io/ioutil"
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
	name1, name2 := args[0], args[1]

	// why are args (text1&2) handled differently here?
	text1 := ts.ReadFile(name1)
	absName2 := ts.MkAbs(name2)
	data, err := ioutil.ReadFile(absName2)
	ts.Check(err)
	text2 := string(data)
	if env {
		text2 = ts.expand(text2)
	}
	if text1 == text2 {
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
	fmt.Println(text1)

	if text1 == text2 {
		return
	}

	ts.Logf("[diff -%s(%d) +%s(%d)]\n%s\n", name1, len(text1), name2, len(text2), diff.Diff(name1, []byte(text1), name2, []byte(text2)))
	ts.Fatalf("%s and %s differ", name1, name2)
}
