package incept

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/google/shlex"
	"github.com/google/uuid"
	"github.com/mattn/go-isatty"

	"dagger.io/dagger/telemetry"
	"github.com/dagger/dagger/dagql/dagui"
	"github.com/dagger/dagger/dagql/idtui"
	"github.com/dagger/dagger/engine/client"
)

type InceptOptions struct {
	// IO configuration
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// Frontend configuration
	Progress string // "auto", "plain", "tty", "dots", "report"
	Workdir  string
	Verbose  int
	Quiet    int
	Silent   bool
	Debug    bool

	// Execution configuration
	WaitDelay          time.Duration
	DisableHostRW      bool
	Interactive        bool
	InteractiveCommand string
	NoExit             bool
}

// Incept mimics the `dagger run ...` command as a function without needing to run dagger.
// The args here are the same you would pass to `dagger run ...`
func Incept(ctx context.Context, args []string, options *InceptOptions) error {
	if options == nil {
		options = &InceptOptions{}
	}

	// 1. Initialize Globals (from main.go and engine.go)
	// We need to set these because withEngine relies on them.

	// HMMM, do we need one of those things that allows us to duplicate (fan out) i/o so we can also print after? (or save to file) (maybe even different frontends?!)

	// IO Globals
	if options.Stdout != nil {
		stdout = options.Stdout
	} else {
		stdout = os.Stdout
	}
	if options.Stderr != nil {
		stderr = options.Stderr
	} else {
		stderr = os.Stderr
	}
	if options.Stdin != nil {
		stdin = options.Stdin
	} else {
		stdin = os.Stdin
	}

	stdoutIsTTY = isFileTerminal(stdout)
	stderrIsTTY = isFileTerminal(stderr)
	hasTTY = stdoutIsTTY || stderrIsTTY

	// Workdir Global
	if options.Workdir != "" {
		workdir = options.Workdir
	} else {
		// Default normalization
		var err error
		workdir, err = NormalizeWorkdir("")
		if err != nil {
			return err
		}
	}

	// Flags Globals
	verbose = options.Verbose
	quiet = options.Quiet
	silent = options.Silent
	debugFlag = options.Debug
	disableHostRW = options.DisableHostRW
	interactive = options.Interactive

	if options.WaitDelay != 0 {
		waitDelay = options.WaitDelay
	} else {
		waitDelay = 10 * time.Second
	}

	if options.InteractiveCommand != "" {
		interactiveCommand = options.InteractiveCommand
	} else {
		interactiveCommand = "/bin/sh"
	}

	var err error
	interactiveCommandParsed, err = shlex.Split(interactiveCommand)
	if err != nil {
		return fmt.Errorf("cannot parse interactive command: %w", err)
	}

	// 2. Setup Frontend and Options (dagui.FrontendOpts)

	// Reset global opts to zero-value-ish defaults, then apply overrides
	opts = dagui.FrontendOpts{}
	opts.Verbosity += dagui.ShowCompletedVerbosity // keep progress by default
	opts.Verbosity += options.Verbose
	opts.Verbosity -= options.Quiet
	opts.Silent = options.Silent
	opts.Debug = options.Debug
	opts.NoExit = options.NoExit
	// opts.RevealNoisySpans = ... (not exposed in options struct yet, defaulting to false)
	// opts.ExpandCompleted = ...

	// Progress
	progress = options.Progress
	if progress == "" {
		progress = "auto"
	}
	if progress == "auto" {
		if hasTTY {
			progress = "tty"
		} else {
			progress = "plain"
		}
	}
	if silent {
		progress = "plain"
	}

	switch progress {
	case "plain":
		Frontend = idtui.NewPlain(stderr)
	case "tty":
		if !hasTTY {
			Frontend = idtui.NewPlain(stderr)
		} else {
			Frontend = idtui.NewPretty(stderr)
		}
	case "dots":
		Frontend = idtui.NewDots(stderr)
	case "report":
		Frontend = idtui.NewReporter(stderr)
	default:
		return fmt.Errorf("unknown progress type %q", progress)
	}

	// 3. Run Logic (similar to run.go)

	u, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generate uuid: %w", err)
	}
	sessionToken := u.String()

	otelEnv, err := setupTelemetryProxy(ctx)
	if err != nil {
		return fmt.Errorf("setup telemetry proxy: %w", err)
	}

	// withEngine takes care of connection, starting the engine, etc.
	return withEngine(ctx, client.Params{
		SecretToken: sessionToken,
	}, func(ctx context.Context, engineClient *client.Client) error {
		sessionL, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return fmt.Errorf("session listen: %w", err)
		}
		defer sessionL.Close()

		env := os.Environ()
		sessionPort := fmt.Sprintf("%d", sessionL.Addr().(*net.TCPAddr).Port)
		env = append(env, "DAGGER_SESSION_PORT="+sessionPort)
		env = append(env, "DAGGER_SESSION_TOKEN="+sessionToken)
		env = append(env, telemetry.PropagationEnv(ctx)...)
		env = append(env, otelEnv...)

		// Prepare subcommand
		if len(args) == 0 {
			return fmt.Errorf("no command specified")
		}
		subCmd := exec.CommandContext(ctx, args[0], args[1:]...)
		subCmd.Env = env
		subCmd.Stdin = stdin
		if !silent {
			stdio := telemetry.SpanStdio(ctx, InstrumentationLibrary)
			if stdoutIsTTY {
				subCmd.Stdout = stdio.Stdout
			} else {
				subCmd.Stdout = stdout
			}
			if stderrIsTTY {
				subCmd.Stderr = stdio.Stderr
			} else {
				subCmd.Stderr = stderr
			}
		} else {
			subCmd.Stdout = stdout
			subCmd.Stderr = stderr
		}

		ensureChildProcessesAreKilled(subCmd)

		// Start Session Server
		srv := &http.Server{
			Handler: engineClient,
			BaseContext: func(listener net.Listener) context.Context {
				return ctx
			},
		}

		go srv.Serve(sessionL)

		// Run Command
		err = subCmd.Run()
		return err
	})
}

func isFileTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return isatty.IsTerminal(f.Fd())
	}
	return false
}
