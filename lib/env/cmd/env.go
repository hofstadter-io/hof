package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env/dag"
	"golang.org/x/sync/errgroup"
)

func Env(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole) error {
	args, cueargs := splitArgs(args)
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
	d, _ := dag.NewClient(R.Ctx, R.DagClient)

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

					seqCtx, seqSpan := dagger.Tracer().Start(R.Ctx, fmt.Sprintf("env.%s.step.%d", e.Hof.Env.Name, s1))

					g, ctx := errgroup.WithContext(seqCtx)
					if eflags.Parallel > 0 {
						g.SetLimit(eflags.Parallel)
					}

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
							seqSpan.End()
							return err
						}
						// fmt.Printf("  [%d/%d][%d/%d]: %s", s1+1, len(seqSteps), s2+1, len(seqStep), k.Name)

						//
						// Phase 1 - eval cue and assemble dagger pointers
						//

						var file *dagger.File
						var allowParent bool
						var dir *dagger.Directory
						var wipe bool
						var dest string

						switch k.Kind {
						case "#exportFile":
							_file, cfg, _err := d.HashExportFile(parStep)
							if cfg != nil {
								dest = cfg.Path
							}
							allowParent = cfg.AllowParentDirPath
							file, err = _file, _err
						case "#exportDir":
							_dir, cfg, _err := d.HashExportDir(parStep)
							if cfg != nil {
								dest = cfg.Path
								wipe = cfg.Wipe
							}
							dir, err = _dir, _err
						case "#exportImage":
							_c, cfg, _err := d.HashExportImage(parStep)
							if cfg != nil {
								dest = cfg.Name
							}
							c, err = _c, _err
						case "#exportImageFile":
							_c, cfg, _err := d.HashExportImageFile(parStep)
							if cfg != nil {
								dest = cfg.Path
							}
							c, err = _c, _err
						case "#publishImage":
							_c, cfg, _err := d.HashPublishImage(parStep)
							if cfg != nil {
								dest = cfg.Name
							}
							c, err = _c, _err

						case "#file":
							file, dest, err = d.File(parStep, eflags.NoCache)
						case "#dir":
							dir, dest, err = d.Dir(parStep, eflags.NoCache)

						case "#container":
							c, err = d.HashContainer(parStep)
						case "#dockerBuild":
							c, err = d.HashDockerBuild(parStep)

						default:
							seqSpan.End()
							return fmt.Errorf("unsupported cmd target(%s): %v", k.Kind, parStep)
						}

						//
						// Phase 2 - synchronize dagger, in parallel
						//

						// we should control parallelism here

						// only parallel the dagger work
						g.Go(func() error {
							parCtx, parSpan := dagger.Tracer().Start(ctx, fmt.Sprintf("step.%s", k.Name))
							defer parSpan.End()

							// if we already have an error, just return it for collection
							if err != nil {
								return err
							}

							start := time.Now()

							// HMMM, this decides what we do

							switch k.Kind {
							// todo, we need to split these across here (cue eval) & below (dag sync)
							case "#exportFile":
								_, err = file.Export(parCtx, dest, dagger.FileExportOpts{
									AllowParentDirPath: allowParent,
								})
							case "#exportDir":
								_, err = dir.Export(parCtx, dest, dagger.DirectoryExportOpts{
									Wipe: wipe,
								})
							case "#exportImage":
								err = c.ExportImage(parCtx, dest, dagger.ContainerExportImageOpts{})
							case "#exportImageFile":
								_, err = c.Export(parCtx, dest, dagger.ContainerExportOpts{})

							case "#publish":
								c, err = c.Sync(parCtx)
							case "#container":
								c, err = c.Sync(parCtx)
							case "#hostImage":
								c, err = c.Sync(parCtx)
							}

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
					seqSpan.End()
					if err != nil {
						return fmt.Errorf("while executing parallel tasks(%s.%s.%d): %w", e.Hof.Env.Name, t, s1, err)
					}

				} // end loop over task-seq-step
			} // end loop over tasks
		} // end loop over envs
	} // end loop over args (cmds)

	return nil
}
