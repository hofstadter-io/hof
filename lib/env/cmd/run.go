package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"dagger.io/dagger"
	"golang.org/x/sync/errgroup"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/env/dag"
)

func runnable(e *env.Env) bool {
	accepting := []string{"container", "hostImage", "dockerBuild", "service"}
	_, kind, _ := extractMeta(e)
	// only publish containers right now
	if slices.Contains(accepting, kind) {
		return true
	}
	return false
}

func Run(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole, cflags flags.Env__RunFlagpole) error {
	// some quick setup and early filtering
	R, matches, err := commonStart(args, rflags, eflags, runnable)
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

	buildCtx, buildCancel := context.WithCancel(R.Ctx)
	defer buildCancel()

	buildCtx, buildSpan := dagger.Tracer().Start(buildCtx, "hof env run")
	defer buildSpan.End()

	var serviceMatches []*env.Env
	var containerMatches []*env.Env

	for _, e := range matches {
		_, kind, _ := extractMeta(e)
		if kind == "service" {
			serviceMatches = append(serviceMatches, e)
		} else {
			containerMatches = append(containerMatches, e)
		}
	}

	var g *errgroup.Group
	var groupCtx context.Context

	if len(serviceMatches) > 0 {
		fmt.Println("starting:")
		if eflags.FailFast {
			g, groupCtx = errgroup.WithContext(buildCtx)
		} else {
			g = new(errgroup.Group)
			groupCtx = buildCtx
		}

		if eflags.Parallel > 0 {
			g.SetLimit(eflags.Parallel)
		}

		for _, e := range serviceMatches {
			e := e
			g.Go(func() error {
				name, kind, _ := extractMeta(e)
				matchCtx, matchSpan := dagger.Tracer().Start(groupCtx, fmt.Sprintf("starting[%s]: (%s)", name, kind))
				defer matchSpan.End()

				if groupCtx.Err() != nil {
					fmt.Printf("ABORT: %s (%s)\n", name, kind)
					return groupCtx.Err()
				}

				fmt.Printf("START: %s (%s)\n", name, kind)
				start := time.Now()

				s, cfg, err := d.Service(e.Value, eflags.NoCache)
				if err != nil {
					fmt.Printf("ERROR: %s (%s) %v\n", name, kind, err)
					return err
				}

				ports := []dagger.PortForward{}
				for _, p := range cfg.Ports {
					if p.Frontend == 0 {
						p.Frontend = p.Backend
					}
					ports = append(ports, dagger.PortForward{
						Backend:  p.Backend,
						Frontend: p.Frontend,
						Protocol: dagger.NetworkProtocol(strings.ToUpper(p.Protocol)),
					})
				}

				s = R.DAG.Host().Tunnel(s, dagger.HostTunnelOpts{
					Ports: ports,
				})
				s, err = s.Start(matchCtx)
				if err != nil {
					if errors.Is(err, context.Canceled) || isDaggerQueryError(err) {
						fmt.Printf("ABORT: %s (%s)\n", name, kind)
					} else {
						fmt.Printf("ERROR: %s (%s) %v\n", name, kind, err)
					}
					return err
				}

				fmt.Printf(" DONE: %s (%s) (%v)\n", name, kind, time.Since(start).Round(time.Millisecond))

				// wait for context to be canceled before stopping
				<-groupCtx.Done()

				fmt.Printf(" STOP: %s (%s)\n", name, kind)
				// create a new context for stopping as groupCtx is already canceled
				stopCtx, stopCancel := context.WithTimeout(context.Background(), 45*time.Second)
				defer stopCancel()

				stopWaitCtx, stopWaitCancel := context.WithTimeout(stopCtx, 30*time.Second)
				_, err = s.Stop(stopWaitCtx)
				stopWaitCancel()

				if err != nil {
					fmt.Printf(" KILL: %s (%s) %v\n", name, kind, err)
					_, err = s.Stop(stopCtx, dagger.ServiceStopOpts{Kill: true})
					if err != nil {
						fmt.Printf("ERROR: %s (%s) %v\n", name, kind, err)
					}
				}

				return nil
			})
		}
	}

	if len(containerMatches) > 0 {
		for i, e := range containerMatches {
			isLast := i == len(containerMatches)-1

			c, err := d.Container(e.Value, eflags.NoCache)
			if err != nil {
				return err
			}

			if isLast {
				var cmd []string
				if cflags.Command != "" {
					cmd = strings.Fields(cflags.Command)
				}

				_, err = c.Terminal(dagger.ContainerTerminalOpts{
					Cmd:                           cmd,
					ExperimentalPrivilegedNesting: eflags.Unsafe,
					InsecureRootCapabilities:      eflags.Unsafe,
				}).Sync(buildCtx)
				if err != nil {
					return err
				}
			} else {
				_, err = c.Sync(buildCtx)
				if err != nil {
					return err
				}
			}
		}
	} else if len(serviceMatches) > 0 {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
	}

	// cancel services and wait
	buildCancel()
	if len(serviceMatches) > 0 {
		return g.Wait()
	}
	err = d.StopWatchers()
	if err != nil {
		return err
	}

	return nil
}
