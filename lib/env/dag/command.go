package dag

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"golang.org/x/sync/errgroup"
)

type hashCmdConfig struct {
	Kind string `json:"$kind"`
	Name string `json:"name"`

	// how do we order these? ... hear me out ... @z(int), same values have no order guarantees
	// gets all mixed up between CUE topo & go maps, but lists aren't fun...
	Tasks  map[string]cue.Value `json:"tasks"`
	Hooks  map[string]cue.Value `json:"hooks"`
	Config map[string]any       `json:"config"`
}

func (d *Dag) DecodeHashCmd(step cue.Value) (*hashCmdConfig, error) {
	var cfg hashCmdConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding hashCmd: %w", err)
	}

	return &cfg, nil
}

type hashTaskConfig struct {
	Kind string `json:"$kind"`
	Name string `json:"name"`

	Steps []cue.Value          `json:"steps"`
	Hooks map[string]cue.Value `json:"hooks"`

	Parallel int            `json:"parallel"`
	Config   map[string]any `json:"config"`
}

func (d *Dag) DecodeHashTask(step cue.Value) (*hashTaskConfig, error) {
	var cfg hashTaskConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding hashTask: %w", err)
	}

	return &cfg, nil
}

type RunTaskOpts struct {
	EnvName  string
	Parallel int
	FailFast bool
	NoCache  bool
}

func (d *Dag) RunTask(ctx context.Context, taskCfg *hashTaskConfig, opts RunTaskOpts) error {
	steps := taskCfg.Steps
	parallel := taskCfg.Parallel
	// override if set in opts (cli flags)
	if opts.Parallel > 0 {
		parallel = opts.Parallel
	}

	// fmt.Printf("      [%d/%d]\n", s1, len(seqSteps))

	taskCtx, taskSpan := dagger.Tracer().Start(ctx, fmt.Sprintf("env.%s.task.%s", opts.EnvName, taskCfg.Name))
	defer taskSpan.End()

	var g *errgroup.Group
	var gCtx context.Context

	if parallel > 0 {
		if opts.FailFast {
			g, gCtx = errgroup.WithContext(taskCtx)
		} else {
			g = new(errgroup.Group)
			gCtx = taskCtx
		}
		g.SetLimit(parallel)
	} else {
		gCtx = taskCtx
	}

	// Helper to process a list of steps recursively
	var processSteps func(steps []cue.Value) error
	processSteps = func(steps []cue.Value) error {
		for i, step := range steps {

			// Handle nested lists (infinite nesting)
			if step.IncompleteKind() == cue.ListKind {
				iter, _ := step.List()
				var substeps []cue.Value
				for iter.Next() {
					substeps = append(substeps, iter.Value())
				}
				if err := processSteps(substeps); err != nil {
					return err
				}
				continue
			}

			// CUE is not concurrency safe yet
			// fmt.Printf("         [%d/%d]\n", s2, len(seqStep))
			type brief struct {
				Kind string `json:"$kind"`
				Name string `json:"name"`
			}
			var k brief
			var err error

			// Resolve #shouldi
			step, err = d.ResolveShouldi(step)
			if err != nil {
				return err
			}
			if !step.Exists() {
				continue
			}

			err = step.Decode(&k)
			if err != nil {
				return err
			}
			// fmt.Printf("  [%d/%d][%d/%d]: %s", s1+1, len(seqSteps), s2+1, len(seqStep), k.Name)

			//
			// Phase 1 - eval cue and assemble dagger pointers
			//

			var c *dagger.Container
			var svc *dagger.Service
			var file *dagger.File
			var allowParent bool
			var dir *dagger.Directory
			var wipe bool
			var dest string

			switch k.Kind {
			case "#exportFile":
				_file, cfg, _err := d.HashExportFile(step, opts.NoCache)
				if cfg != nil {
					dest = cfg.Path
				}
				allowParent = cfg.AllowParentDirPath
				file, err = _file, _err
			case "#exportDir":
				_dir, cfg, _err := d.HashExportDir(step, opts.NoCache)
				if cfg != nil {
					dest = cfg.Path
					wipe = cfg.Wipe
				}
				dir, err = _dir, _err
			case "#exportImage":
				_c, cfg, _err := d.HashExportImage(step, opts.NoCache)
				if cfg != nil {
					dest = cfg.Name
				}
				c, err = _c, _err
			case "#exportImageFile":
				_c, cfg, _err := d.HashExportImageFile(step, opts.NoCache)
				if cfg != nil {
					dest = cfg.Path
				}
				c, err = _c, _err
			case "#publishImage":
				_c, cfg, _err := d.HashPublishImage(step, opts.NoCache)
				if cfg != nil {
					dest = cfg.Name
				}
				c, err = _c, _err

			case "#hostExec":
				// no-op in Phase 1, we run in Phase 2

			case "#file":
				file, dest, err = d.File(step, opts.NoCache)
			case "#dir":
				dir, dest, err = d.Dir(step, opts.NoCache)
			case "#rootfs":
				dir, dest, err = d.Dir(step, opts.NoCache)

			case "#container":
				c, err = d.HashContainer(step, opts.NoCache)
			case "#service":
				svc, _, err = d.HashService(step, opts.NoCache)
			case "#dockerBuild":
				c, err = d.HashDockerBuild(step, opts.NoCache)

			default:
				return fmt.Errorf("unsupported cmd target(%s): %v", k.Kind, step)
			}

			if err != nil {
				return fmt.Errorf("in step %s: %w", k.Name, err)
			}

			//
			// Phase 2 - synchronize dagger, in parallel
			//

			// execFunc does the actual work
			execFunc := func(ctx context.Context) error {
				ctx, span := dagger.Tracer().Start(ctx, fmt.Sprintf("step.%s", k.Name))
				defer span.End()

				// if we already have an error, just return it for collection
				if ctx.Err() != nil {
					fmt.Printf("ABORT: %s\n", k.Name)
					return ctx.Err()
				}

				fmt.Printf("START: %s\n", k.Name)
				start := time.Now()

				// HMMM, this decides what we do

				switch k.Kind {
				// todo, we need to split these across here (cue eval) & below (dag sync)
				case "#hostExec":
					err = d.HashHostExec(step)
				case "#exportFile":
					_, err = file.Export(ctx, dest, dagger.FileExportOpts{
						AllowParentDirPath: allowParent,
					})
				case "#exportDir":
					_, err = dir.Export(ctx, dest, dagger.DirectoryExportOpts{
						Wipe: wipe,
					})
				case "#exportImage":
					err = c.ExportImage(ctx, dest, dagger.ContainerExportImageOpts{})
				case "#exportImageFile":
					_, err = c.Export(ctx, dest, dagger.ContainerExportOpts{})

				case "#publish":
					c, err = c.Sync(ctx)
				case "#container":
					c, err = c.Sync(ctx)
				case "#service":
					_, err = svc.Start(ctx)
				case "#hostImage":
					c, err = c.Sync(ctx)
				}

				// TODO, build up or exit, depending on config
				str := fmt.Sprintf("%s.[%d/%d]", k.Name, i+1, len(steps))
				if err != nil {
					if errors.Is(err, context.Canceled) || isDaggerQueryError(err) {
						fmt.Printf("ABORT: %s\n", k.Name)
					} else {
						fmt.Printf("ERROR: %s %v\n", str, err)
					}
					return err
				}
				fmt.Printf(" DONE: %s (%v)\n", str, time.Since(start).Round(time.Millisecond))

				return err
			}

			if parallel > 0 {
				g.Go(func() error {
					return execFunc(gCtx)
				})
			} else {
				if err := execFunc(gCtx); err != nil {
					return err
				}
			}

		} // end loop over task-seq-step
		return nil
	}

	if err := processSteps(steps); err != nil {
		return err
	}

	if parallel > 0 {
		return g.Wait()
	}

	return nil
}

func isDaggerQueryError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	if !strings.HasPrefix(s, `Post "`) {
		return false
	}
	// find the first quote after Post "
	endIdx := strings.Index(s[6:], `"`)
	if endIdx == -1 {
		return false
	}
	uriStr := s[6 : 6+endIdx]
	u, err := url.Parse(uriStr)
	if err != nil {
		return false
	}
	return u.Path == "/query"
}
