package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env/dag"
	"github.com/hofstadter-io/hof/lib/env/incept"
	"golang.org/x/sync/errgroup"
)

func Env(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole) error {
	args, cueargs := splitArgs(args)
	R, err := prepRuntime(cueargs, rflags)
	if err != nil {
		return err
	}

	// incept if we are not in dagger
	dst := os.Getenv("DAGGER_SESSION_TOKEN")
	if dst == "" {
		// Run incept
		err := incept.Incept(context.Background(), os.Args, &incept.InceptOptions{
			Progress:    eflags.Progress,
			Interactive: eflags.OnFailure,
			NoExit:      eflags.NoExit,
			Stdout:      os.Stdout,
			Stderr:      os.Stderr,
			Stdin:       os.Stdin,
		})
		if err != nil {
			return err
		}

		return nil
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	if err != nil {
		return fmt.Errorf("while connecting to dagger in build: %w", err)
	}
	d, _ := dag.NewClient(ctx, client)

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
				// fmt.Printf("   [%d/%d]: %s\n", t1, len(cmdCfg.Tasks), t)
				// fmt.Printf("  %s:\n", t)
				// fmt.Println(taskVal)
				taskCfg, err := d.DecodeHashTask(taskVal)
				if err != nil {
					return err
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

				var c *dagger.Container
				seqSteps := taskCfg.Steps
				for s1, seqStep := range seqSteps {
					// fmt.Printf("      [%d/%d]\n", s1, len(seqSteps))

					g, ctx := errgroup.WithContext(context.Background())

					for s2, parStep := range seqStep {

						// CUE is not concurrency safe yet
						// fmt.Printf("         [%d/%d]\n", s2, len(seqStep))
						type brief struct {
							Kind string `json:"$kind"`
							Name string `json:"name"`
						}
						var k brief
						err := parStep.Decode(&k)
						if err != nil {
							return err
						}
						// fmt.Printf("  [%d/%d][%d/%d]: %s", s1+1, len(seqSteps), s2+1, len(seqStep), k.Name)

						// go do() something with parStep
						switch k.Kind {
						case "#container":
							c, err = d.HashContainer(parStep)
						case "#hostImage":
							c, err = d.HashHostImage(parStep)
						default:
							return fmt.Errorf("unsupported cmd target(%s): %v", k.Kind, parStep)
						}

						// only parallel the dagger work
						g.Go(func() error {

							start := time.Now()

							// TODO, build up or exit, depending on config
							if err != nil {
								return err
							}

							c, err = c.Sync(ctx)

							// TODO, build up or exit, depending on config
							str := fmt.Sprintf("%s.[%d/%d]", k.Name, s1+1, s2+1)
							outcome := "ok"
							if err != nil {
								outcome = "err"
							}
							fmt.Printf("%-32s  %-3s  %s\n", str, outcome, time.Since(start))

							return err
						}) // end of goroutine
					} // end loop spawning parallel containers for the task-step

					err = g.Wait()
					if err != nil {
						return fmt.Errorf("while executing parallel tasks(%s.%s.%d): %w", e.Hof.Env.Name, t, s1, err)
					}

				} // end loop over task-seq-step

				// shell at the end of a task
				// if eflags.Shell {
				// 	c = c.Terminal()
				// 	c, err = c.Sync(ctx)
				// 	if err != nil {
				// 		return err
				// 	}
				// }

			} // end loop over tasks
		} // end loop over envs
	} // end loop over args (cmds)

	return nil
}
