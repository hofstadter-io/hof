package runtime

import (
	"io/ioutil"
	"strings"
)

// env displays or adds to the environment.
func (ts *Script) CmdEnv(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? env")
	}
	if len(args) == 0 {
		printed := make(map[string]bool) // env list can have duplicates; only print effective value (from envMap) once
		for _, kv := range ts.env {
			k := envvarname(kv[:strings.Index(kv, "=")])
			if !printed[k] {
				printed[k] = true
				ts.Logf("%s=%s\n", k, ts.envMap[k])
			}
		}
		return
	}
	for _, env := range args {
		i := strings.Index(env, "=")
		if i < 0 {
			// Display value instead of setting it.
			ts.Logf("%s=%s\n", env, ts.Getenv(env))
			continue
		}
		k, v := env[:i], env[i+1:]
		if v[0] == '@' {
			fname := v[1:] // for error messages
			if fname == "stdout" {
				v = ts.stdout
			} else if fname == "stderr" {
				v = ts.stderr
			} else {
				data, err := ioutil.ReadFile(ts.MkAbs(fname))
				ts.Check(err)
				v = string(data)
			}
		}
		ts.Setenv(k,v)
	}
}

// sub any env vars in a string or file
func (ts *Script) CmdEnvsub(neg int, args []string) {
	if neg != 0 {
		// It would be strange to say "this file can have any content except this precise byte sequence".
		ts.Fatalf("unsupported: !? cmp")
	}
	if len(args) != 1 {
		ts.Fatalf("usage: envsub string/@file")
	}

	v := args[0]

	if v[0] == '@' {
		fname := v[1:] // for error messages
		if fname == "stdout" {
			v = ts.stdout
		} else if fname == "stderr" {
			v = ts.stderr
		} else {
			data, err := ioutil.ReadFile(ts.MkAbs(fname))
			ts.Check(err)
			v = string(data)
		}
	}

	subd := ts.expand(v)
	ts.stdout = subd
}

