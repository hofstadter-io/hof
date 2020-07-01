package runtime

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// scriptMatch implements both stdout and stderr.
func scriptMatch(ts *Script, neg int, args []string, text, name string) {
	n := 0
	if len(args) >= 1 && strings.HasPrefix(args[0], "-count=") {
		if neg != 0 {
			ts.Fatalf("cannot use -count= with negated match")
		}
		var err error
		n, err = strconv.Atoi(args[0][len("-count="):])
		if err != nil {
			ts.Fatalf("bad -count=: %v", err)
		}
		if n < 1 {
			ts.Fatalf("bad -count=: must be at least 1")
		}
		args = args[1:]
	}

	isRegexp := name == "regexp"
	isGrep := name == "grep"
	isSed := name == "sed"

	extraUsage := ""
	want := 1
	if isRegexp || isGrep {
		extraUsage = " file"
		want = 2
	}
	if isSed {
		extraUsage = " replace file"
		want = 3
	}
	if len(args) != want {
		ts.Fatalf("usage: %s [-count=N] 'pattern'%s", name, extraUsage)
	}

	pattern := args[0]
	switch pattern {
	case "stdout":
		pattern = ts.stdout
	case "stderr":
		pattern = ts.stderr

	default:
		if pattern[0] == '@' {
			fname := pattern[1:] // for error messages
			data, err := ioutil.ReadFile(ts.MkAbs(fname))
			ts.Check(err)
			pattern = string(data)
		}
	}
	re, err := regexp.Compile(`(?m)` + pattern)
	ts.Check(err)


	if isRegexp || isGrep {
		content := args[1]
		switch  content {
		case "stdout", "$WORK/stdout":
			text = ts.stdout
		case "stderr", "$WORK/stderr":
			text = ts.stderr

		default:
			name = args[1] // for error messages
			data, err := ioutil.ReadFile(ts.MkAbs(args[1]))
			ts.Check(err)
			text = string(data)
		}
	}
	replace := ""
	if isSed {
		replace = args[1]
		switch  replace {
		case "stdout", "$WORK/stdout":
			text = ts.stdout
		case "stderr", "$WORK/stderr":
			text = ts.stderr

		default:
			if replace[0] == '@' {
				fname := replace[1:] // for error messages
				data, err := ioutil.ReadFile(ts.MkAbs(fname))
				ts.Check(err)
				replace = string(data)
			}
		}
		content := args[2]
		switch  content {
		case "stdout", "$WORK/stdout":
			text = ts.stdout
		case "stderr", "$WORK/stderr":
			text = ts.stderr

		default:
			if content[0] == '@' {
				fname := content[1:] // for error messages
				data, err := ioutil.ReadFile(ts.MkAbs(fname))
				ts.Check(err)
				content = string(data)
			}
		}
	}

	if neg > 0 {
		if re.MatchString(text) {
			if isGrep {
				ts.Logf("[%s]\n%s\n", name, text)
			}
			ts.Fatalf("unexpected match for %#q found in %s: %s", pattern, name, re.FindString(text))
		}

		if isGrep {
			c := -1
			if n > 0 {
				c = n
			}
			matches := re.FindAllString(text, c)
			if c > 0 && len(matches) > c {
				matches = matches[:c]
			}
			ts.stdout = strings.Join(matches, "\n")
		}
		if isSed {
			ts.stdout = re.ReplaceAllString(text, replace)
		}
	} else {
		if isGrep || isSed {
			ts.Fatalf("%s does not support status checking", name)
		}
		if !re.MatchString(text) {
			if isGrep {
				ts.Logf("[%s]\n%s\n", name, text)
			}
			ts.Fatalf("no match for %#q found in %s", pattern, name)
		}
		if n > 0 {
			count := len(re.FindAllString(text, -1))
			if count != n {
				if isGrep {
					ts.Logf("[%s]\n%s\n", name, text)
				}
				ts.Fatalf("have %d matches for %#q, want %d", count, pattern, n)
			}
		}
	}
}

