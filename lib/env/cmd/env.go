package cmd

import (
	"regexp"
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/env/dag"
)

func Env(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole) error {

	args, cueargs := cuetils.PercentSplitArgs(args)
	R, err := prepRuntime(cueargs, rflags)
	if err != nil {
		return err
	}

	// incept if we are not in dagger
	incepted, err := daggerInceptFlags(rflags, eflags)
	if incepted {
		return err
	}

	// setup dagger & cue->dagger engine
	err = R.DaggerInit()
	d, _ := dag.NewClient(R.Ctx, R.DAG)

	var cmdArg, taskArg string
	for _, arg := range args {
		cmdArg = arg
		if strings.Contains(cmdArg, "/") {
			parts := strings.Split(cmdArg, "/")
			cmdArg = parts[0]
			taskArg = strings.Join(parts[1:], "/")
		}

		// fmt.Println("args:", cmdArg, taskArg)

		for _, e := range R.Envs {
			// fmt.Println("env:", e.Hof.Env.Kind, e.Hof.Env.Name)
			if e.Hof.Env.Kind != "cmd" || !(cmdArg == "" || cmdArg == e.Hof.Env.Name) {
				continue
			}
			// fmt.Printf("%s:\n", cmdArg)
			// fmt.Println(e.Value)

			cmdCfg, err := d.DecodeHashCmd(e.Value)
			if err != nil {
				return err
			}

			t1 := 0
			for t, taskVal := range cmdCfg.Tasks {
				t1++
				taskCfg, err := d.DecodeHashTask(taskVal)
				if err != nil {
					return err
				}
				if taskCfg.Name == "" {
					taskCfg.Name = t
				}

				if taskArg != "" {
					re, err := regexp.Compile(taskArg)
					if err != nil {
						return err
					}
					if !re.MatchString(t) {
						continue
					}
				}

				err = d.RunTask(R.Ctx, taskCfg, dag.RunTaskOpts{
					EnvName:  e.Hof.Env.Name,
					Parallel: eflags.Parallel,
					FailFast: eflags.FailFast,
					NoCache:  eflags.NoCache,
				})
				if err != nil {
					return err
				}
			} // end loop over tasks
		} // end loop over envs
	} // end loop over args (cmds)

	return nil
}
