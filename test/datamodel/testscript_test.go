package datamodel_test

import (
	"os"
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func envSetup(env *runtime.Env) error {
	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")
	env.Vars = append(env.Vars, "GITHUB_TOKEN=" + os.Getenv("GITHUB_TOKEN"))
	env.Vars = append(env.Vars, "HOF_FMT_VERSION=" + os.Getenv("HOF_FMT_VERSION"))
	return nil
}

func TestDatamodel(t *testing.T) {
	yagu.Mkdir(".workdir/testdata")
	runtime.Run(t, runtime.Params{
		Dir:         "./testdata",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/testdata",
		Setup:       envSetup,
	})
}

//func TestBugcases(t *testing.T) {
	//yagu.Mkdir(".workdir/bugcases")
	//runtime.Run(t, runtime.Params{
		//Dir:         "./bugcases",
		//Glob:        "*.txt",
		//WorkdirRoot: ".workdir/bugcases",
		//Setup:       envSetup,
	//})
//}

