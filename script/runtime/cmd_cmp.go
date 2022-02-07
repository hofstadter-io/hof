package runtime

import (
	"io/ioutil"
  "strings"

	"github.com/hofstadter-io/hof/lib/gotils/intern/textutil"
)

// cmp compares two files.
func (ts *Script) CmdCmp(neg int, args []string) {
	if neg != 0 {
		// It would be strange to say "this file can have any content except this precise byte sequence".
		ts.Fatalf("unsupported: !? cmp")
	}
	if len(args) != 2 {
		ts.Fatalf("usage: cmp file1 file2")
	}

	ts.doCmdCmp(args, false)
}

// cmpenv compares two files with environment variable substitution.
func (ts *Script) CmdCmpenv(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? cmpenv")
	}
	if len(args) != 2 {
		ts.Fatalf("usage: cmpenv file1 file2")
	}
	ts.doCmdCmp(args, true)
}

func (ts *Script) doCmdCmp(args []string, env bool) {
	name1, name2 := args[0], args[1]
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

  text3 := strings.TrimSpace(text1)
  text4 := strings.TrimSpace(text2)
  s := ""
  if text3 == text4 {
    s = "(in leading | trailing whitespace)\n"
  }
	ts.Logf("[diff -%s +%s]\n%s%s\n", name1, name2, s, textutil.Diff(text1, text2))
	ts.Fatalf("%s and %s differ", name1, name2)
}

