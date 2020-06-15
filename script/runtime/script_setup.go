package runtime

import (
	"io/ioutil"
	"os"
	"path/filepath"
	goruntime "runtime"
	"strings"

	"github.com/hofstadter-io/hof/lib/gotils/txtar"
)

// setup sets up the test execution temporary directory and environment.
// It returns the comment section of the txtar archive.
func (ts *Script) setupTest() string {
	ts.workdir = filepath.Join(ts.testTempDir, "script-"+ts.name)
	ts.Check(os.MkdirAll(filepath.Join(ts.workdir, "tmp"), 0777))
	env := &Env{
		Vars: []string{
			"WORK=" + ts.workdir, // must be first for ts.abbrev
			"PATH=" + os.Getenv("PATH"),
			"USER=" + os.Getenv("USER"),
			homeEnvName() + "=/no-home",
			tempEnvName() + "=" + filepath.Join(ts.workdir, "tmp"),
			"devnull=" + os.DevNull,
			"/=" + string(os.PathSeparator),
			":=" + string(os.PathListSeparator),
		},
		WorkDir: ts.workdir,
		Values:  make(map[interface{}]interface{}),
		Cd:      ts.workdir,
		ts:      ts,
	}
	// Must preserve SYSTEMROOT on Windows: https://github.com/golang/go/issues/25513 et al
	if goruntime.GOOS == "windows" {
		env.Vars = append(env.Vars,
			"SYSTEMROOT="+os.Getenv("SYSTEMROOT"),
			"exe=.exe",
		)
	} else {
		env.Vars = append(env.Vars,
			"exe=",
		)
	}
	ts.cd = env.Cd
	// Unpack archive.
	a, err := txtar.ParseFile(ts.file)
	ts.Check(err)
	ts.archive = a
	for _, f := range a.Files {
		name := ts.MkAbs(ts.expand(f.Name))
		ts.scriptFiles[name] = f.Name
		ts.Check(os.MkdirAll(filepath.Dir(name), 0777))
		ts.Check(ioutil.WriteFile(name, f.Data, 0666))
	}
	// Run any user-defined setup.
	if ts.params.Setup != nil {
		ts.Check(ts.params.Setup(env))
	}
	ts.cd = env.Cd
	ts.env = env.Vars
	ts.values = env.Values

	ts.envMap = make(map[string]string)
	for _, kv := range ts.env {
		if i := strings.Index(kv, "="); i >= 0 {
			ts.envMap[envvarname(kv[:i])] = kv[i+1:]
		}
	}
	return string(a.Comment)
}

// setup sets up the test execution temporary directory and environment.
// It returns the comment section of the txtar archive.
func (ts *Script) setupRun() string {
	// ts.workdir = filepath.Join(ts.testTempDir, "script-"+ts.name)
	// ts.Check(os.MkdirAll(filepath.Join(ts.workdir, "tmp"), 0777))

	// expose external ENV here

	env := &Env{
		Vars: os.Environ(),
		WorkDir: ts.workdir,
		Values:  make(map[interface{}]interface{}),
		Cd:      ts.workdir,
		ts:      ts,
	}

	// Must preserve SYSTEMROOT on Windows: https://github.com/golang/go/issues/25513 et al
	if goruntime.GOOS == "windows" {
		env.Vars = append(env.Vars,
			"SYSTEMROOT="+os.Getenv("SYSTEMROOT"),
			"exe=.exe",
		)
	} else {
		env.Vars = append(env.Vars,
			"exe=",
		)
	}

	ts.cd = env.Cd

	// Unpack archive.
	a, err := txtar.ParseFile(ts.file)
	ts.Check(err)
	ts.archive = a
	for _, f := range a.Files {
		name := ts.MkAbs(ts.expand(f.Name))
		ts.scriptFiles[name] = f.Name
		ts.Check(os.MkdirAll(filepath.Dir(name), 0777))
		ts.Check(ioutil.WriteFile(name, f.Data, 0666))
	}

	// Run any user-defined setup.
	if ts.params.Setup != nil {
		ts.Check(ts.params.Setup(env))
	}

	// setup more values on ts from env
	ts.cd = env.Cd
	ts.env = env.Vars
	ts.values = env.Values

	// fmt.Println("EnvArr:", ts.env)

	ts.envMap = make(map[string]string)
	for _, kv := range ts.env {
		if i := strings.Index(kv, "="); i >= 0 {
			ts.envMap[envvarname(kv[:i])] = kv[i+1:]
		}
	}

	return string(a.Comment)
}

