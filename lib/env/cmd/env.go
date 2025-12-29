package cmd

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/env/dag"
	"github.com/hofstadter-io/hof/lib/env/incept"
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
			Interactive: true,
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

	for _, cmd := range args {
		var e *env.Env
		for _, ee := range R.Envs {
			if ee.Hof.Env.Kind == "cmd" && cmd == ee.Hof.Env.Name {
				e = ee
				break
			}
		}

		if e == nil {
			return fmt.Errorf("failed to find cmd %q", cmd)
		}

		fmt.Printf("%s:\n", cmd)
		// fmt.Println(e.Value)

		cmdCfg, err := d.DecodeHashCmd(e.Value)
		if err != nil {
			return err
		}

		t1 := 0
		for t, taskVal := range cmdCfg.Tasks {
			t1++
			// fmt.Printf("   [%d/%d]: %s\n", t1, len(cmdCfg.Tasks), t)
			fmt.Printf("   %s:\n", t)
			// fmt.Println(taskVal)
			taskCfg, err := d.DecodeHashTask(taskVal)
			if err != nil {
				return err
			}

			seqSteps := taskCfg.Steps
			for s1, seqStep := range seqSteps {
				// fmt.Printf("      [%d/%d]\n", s1, len(seqSteps))
				for s2, parStep := range seqStep {
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
					str := fmt.Sprintf("%s.[%d/%d]", k.Name, s1+1, s2+1)
					fmt.Printf("      %-24s", str)

					// go do() something with parStep
					var c *dagger.Container
					switch k.Kind {
					case "#container":
						c, err = d.HashContainer(parStep)
					case "#hostImage":
						c, err = d.HashHostImage(parStep)
					default:
						return fmt.Errorf("unsupported cmd target(%s): %v", k.Kind, parStep)
					}

					// TODO, build up or exit, depending on config
					if err != nil {
						fmt.Printf(" err\n")
						return err
					}

					c, err = c.Sync(ctx)

					// TODO, build up or exit, depending on config
					if err != nil {
						fmt.Printf(" err\n")
						return err
					} else {
						fmt.Printf(" ok\n")
					}

				}
			}
		}

		// i, err = i.Sync(ctx)
		// if err != nil {
		// 	return err
		// }

	}

	return nil
}
