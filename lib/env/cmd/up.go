package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"dagger.io/dagger"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env/dag"
	"github.com/hofstadter-io/hof/lib/env/incept"
)

func Up(args []string, rflags flags.RootPflagpole, cflags flags.EnvPflagpole) error {
	args, cueargs := splitArgs(args)

	// check the runtime first before starting dagger
	R, err := prepRuntime(cueargs, rflags)
	if err != nil {
		return err
	}

	// incept if we are not in dagger
	dst := os.Getenv("DAGGER_SESSION_TOKEN")
	if dst == "" {
		err := incept.Incept(context.Background(), os.Args, &incept.InceptOptions{
			Progress:    cflags.Progress,
			Interactive: cflags.OnFailure,
			NoExit:      cflags.NoExit,
			Stdout:      os.Stdout,
			Stderr:      os.Stderr,
			Stdin:       os.Stdin,
		})
		if err != nil {
			return err
		}

		return nil
	}

	// do normal build stuff

	// this should be on the runtime probable?
	ctx := context.Background()
	os.Setenv("_EXPERIMENTAL_DAGGER_RUNNER_HOST", DAGGER_HOST)

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return fmt.Errorf("while connecting to dagger: %w", err)
	}
	d, _ := dag.NewClient(ctx, client)

	fmt.Println("starting:")
	for _, e := range R.Envs {
		// only building containers right now
		if e.Hof.Env.Kind != "service" {
			continue
		}
		// fmt.Println("-:", e.Hof.Env.Name, e.Hof.Env.Kind)
		// we just try to "build" everything unless there are args
		do := true
		if len(args) > 0 {
			do = false
			for _, a := range args {
				if strings.HasPrefix(e.Hof.Env.Name, a) {
					do = true
					break
				}
			}
		}
		if do {
			fmt.Println(" -", e.Hof.Env.Name)

			s, cfg, err := d.Service(e, cflags.NoCache)
			if err != nil {
				fmt.Println("error:", err)
				return err
			}

			ports := []dagger.PortForward{}
			for _, p := range cfg.Ports {
				ports = append(ports, dagger.PortForward{
					Backend:  p.Backend,
					Frontend: p.Frontend,
					Protocol: dagger.NetworkProtocol(strings.ToUpper(p.Protocol)),
				})
			}

			s = client.Host().Tunnel(s, dagger.HostTunnelOpts{
				Ports: ports,
			})
			s, err = s.Start(ctx)
			if err != nil {
				return err
			}

		}
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return nil
}
