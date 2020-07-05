package runtime

import (
	"os"
	goruntime "runtime"
	"strings"
)

func (RT *Runtime) setupEnv() {
	if RT.envMap == nil {
		RT.envMap = make(map[string]string)
	}

	envs := os.Environ()
	for _, env := range envs {
		i := strings.Index(env, "=")
		k, v := env[:i], env[i+1:]
		RT.envMap[k] = v
	}

	if goruntime.GOOS == "windows" {
		RT.envMap["SYSTEMROOT"] = os.Getenv("SYSTEMROOT")
		RT.envMap["exe"] = ".exe"
	} else {
		RT.envMap["exe"] = ""
	}
}

/*
// setupTest sets up the test execution temporary directory and environment.
// It returns the script content section of the txtar archive.
func (ts *Script) setupTest() string {
	ts.workdir = filepath.Join(ts.testTempDir, "script-"+ts.name)
	ts.Check(os.MkdirAll(filepath.Join(ts.workdir, "tmp"), 0777))
	env := &Env{
		Vars: []string{
			"WORK=" + ts.workdir, // must be first for ts.abbrev
			"PATH=" + os.Getenv("PATH"),
			"USER=" + os.Getenv("USER"),
			homeEnvName() + "=" + ts.workdir + "/home",
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

	return ts.setupFromEnv(env)
}

// setupRun sets up the script execution for working in the current directory.
// the current environment will be exposed to the script
// It returns the script content section of the txtar archive.
func (ts *Script) setupRun() string {

	// expose external ENV here
	env := &Env{
		Vars: os.Environ(),
		WorkDir: ts.workdir,
		Values:  make(map[interface{}]interface{}),
		Cd:      ts.workdir,
		ts:      ts,
	}

	return ts.setupFromEnv(env)
}
*/
