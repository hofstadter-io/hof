package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"

	"dagger.io/dagger"

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
	d, _ := dag.NewClient(R.Ctx, R.DagClient)

	fmt.Println("starting:")
	for _, e := range matches {
		name, kind, _ := extractMeta(e)
		fmt.Printf("  %s (%s)", name, kind)

		s, cfg, err := d.Service(e.Value, eflags.NoCache)
		if err != nil {
			fmt.Println("error:", err)
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

		s = R.DagClient.Host().Tunnel(s, dagger.HostTunnelOpts{
			Ports: ports,
		})
		s, err = s.Start(R.Ctx)
		if err != nil {
			return err
		}

	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return nil
}
