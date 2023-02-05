package mod_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func envSetup(env *runtime.Env) error {
	vars := []string{
		"GITHUB_TOKEN",
		"GITLAB_TOKEN",
		"BITBUCKET_USERNAME",
		"BITBUCKET_PASSWORD",
		"HOFMOD_SSHKEY",
	}
	for _, v := range vars {
		if val := os.Getenv(v); val != "" {
			env.Vars = append(env.Vars, fmt.Sprintf("%s=%s", v, val))
		}
	}
	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")
	return nil
}

func setupWorkdir(dir string) {
	os.RemoveAll(dir)
	yagu.Mkdir(dir)
}

func TestModTests(t *testing.T) {
	d := ".workdir/tests"
	setupWorkdir(d)
	runtime.Run(t, runtime.Params{
		Setup:       envSetup,
		Dir:         "testdata",
		Glob:        "*.txt",
		WorkdirRoot: d,
	})
}

/*
func TestModBugs(t *testing.T) {
	d := ".workdir/bugs"
	setupWorkdir(d)
	runtime.Run(t, runtime.Params{
		Setup:       envSetup,
		Dir:         "testdata/bugs",
		Glob:        "*.txt",
		WorkdirRoot: d,
	})
}
*/

func TestModAuthdApikeysTests(t *testing.T) {
	d := ".workdir/authd/apikeys"
	setupWorkdir(d)
	runtime.Run(t, runtime.Params{
		Setup:       envSetup,
		Dir:         "testdata/authd/apikeys",
		Glob:        "*.txt",
		WorkdirRoot: d,
	})
}

/*
func TestModAuthdSshconfigTests(t *testing.T) {
	d := ".workdir/authd/sshconfig"
	setupWorkdir(d)
	runtime.Run(t, runtime.Params{
		Setup:       envSetup,
		Dir:         "testdata/authd/sshconfig",
		Glob:        "*.txt",
		WorkdirRoot: d,
	})
}
*/

//func TestModAuthdSshkeyTests(t *testing.T) {
//yagu.Mkdir(".workdir/authd/sshkey")
//runtime.Run(t, runtime.Params{
//Setup: envSetup,
//Dir: "testdata/authd/sshkey",
//Glob: "*.txt",
//WorkdirRoot: ".workdir/authd/sshkey",
//})
//}
