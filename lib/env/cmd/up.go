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

func upable(e *env.Env) bool {
	accepting := []string{"service"}
	_, kind, _ := extractMeta(e)
	// only publish containers right now
	if slices.Contains(accepting, kind) {
		return true
	}
	return false
}

func Up(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole) error {
	// some quick setup and early filtering
	R, matches, err := commonStart(args, rflags, eflags, upable)
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

	buildCtx, buildSpan := dagger.Tracer().Start(R.Ctx, "hof env up")
	defer buildSpan.End()

	fmt.Println("starting:")
	var g *errgroup.Group
	var groupCtx context.Context
	if eflags.FailFast {
		g, groupCtx = errgroup.WithContext(buildCtx)
	} else {
		g = new(errgroup.Group)
		groupCtx = buildCtx
	}

	if eflags.Parallel > 0 {
		g.SetLimit(eflags.Parallel)
	}

	for _, e := range matches {
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
			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		return err
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return nil
}
