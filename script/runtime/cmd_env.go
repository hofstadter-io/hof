package runtime

import (
	"io/ioutil"
	"strings"
)

// Setenv sets the value of the environment variable named by the key.
func (ts *Script) Setenv(key, value string) {
	ts.env = append(ts.env, key+"="+value)
	ts.envMap[envvarname(key)] = value
}

// Getenv gets the value of the environment variable named by the key.
func (ts *Script) Getenv(key string) string {
	return ts.envMap[envvarname(key)]
}


// env displays or adds to the environment.
func (ts *Script) CmdEnv(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? env")
	}

	// Print all env vars
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

	// loop over args
	for _, env := range args {

		i := strings.Index(env, "=")

		// if it does not have an '=', then print
		if i < 0 {
			// Display value instead of setting it.
			ts.Logf("%s=%s\n", env, ts.Getenv(env))
			continue
		}

		// else, split and do things
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

		// set the env var
		ts.Setenv(k,v)
	}
}

// env displays or adds to the environment.
func (ts *Script) CmdEnvsub(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? envsub")
	}

	// Print all env vars
	if len(args) == 0 {
		printed := make(map[string]bool) // env list can have duplicates; only print effective value (from envMap) once
		for _, kv := range ts.env {
			k := envvarname(kv[:strings.Index(kv, "=")])
			if !printed[k] {
				printed[k] = true
				subd := ts.expand(ts.envMap[k])
				ts.Logf("%s=%s\n", k, subd)
			}
		}
		return
	}

	// loop over args
	for _, env := range args {

		i := strings.Index(env, "=")

		// if it does not have an '=', then print
		if i < 0 {
			// Display value instead of setting it.
			subd := ts.expand(ts.Getenv(env))
			ts.Logf("%s=%s\n", env, subd)
			continue
		}

		// else, split and do things
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

		// set the env var
		subd := ts.expand(v)
		ts.Setenv(k,subd)
	}
}

