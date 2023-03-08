package runtime

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// expand applies environment variable expansion to the string s.
func (ts *Script) expand(s string) string {
	return os.Expand(s, func(key string) string {
		if key1 := strings.TrimSuffix(key, "@R"); len(key1) != len(key) {
			return regexp.QuoteMeta(ts.Getenv(key1))
		}
		return ts.Getenv(key)
	})
}

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
		parts := strings.Split(env, "=")
		k, v := parts[0], ""
		if len(parts) > 1 {
			v = parts[1]
		}
		if len(v) > 0 && v[0] == '@' {
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
		ts.Setenv(k, v)
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
		ts.Setenv(k, subd)
	}
}

// MkAbs interprets file relative to the test script's current directory
// and returns the corresponding absolute path.
func (ts *Script) MkAbs(file string) string {
	if filepath.IsAbs(file) {
		return file
	}
	return filepath.Join(ts.cd, file)
}

// ReadFile returns the contents of the file with the
// given name, interpreted relative to the test script's
// current directory. It interprets "stdout" and "stderr" to
// mean the standard output or standard error from
// the most recent exec or wait command respectively.
//
// If the file cannot be read, the script fails.
func (ts *Script) ReadFile(file string) string {
	switch file {
	case "stdout":
		return ts.stdout
	case "stderr":
		return ts.stderr
	default:
		file = ts.MkAbs(file)
		data, err := ioutil.ReadFile(file)
		ts.Check(err)
		return string(data)
	}
}

func removeAll(dir string) error {
	// module cache has 0444 directories;
	// make them writable in order to remove content.
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // ignore errors walking in file system
		}
		if info.IsDir() {
			os.Chmod(path, 0777)
		}
		return nil
	})
	return os.RemoveAll(dir)
}
