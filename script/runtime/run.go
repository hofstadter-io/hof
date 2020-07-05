package runtime

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/script/ast"
)

func RunScript(args []string) error {
	params := &Params{
		Mode: Run,
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fs := osfs.New(cwd)

	llvl := "warn"
	if flags.RootPflags.Verbose != "" {
		llvl = flags.RootPflags.Verbose
	}

	config := &ast.Config{
		LogLevel: llvl,
		FS: fs,
	}
	parser := ast.NewParser(config)

	RT := &Runtime{
		params: params,
		parser: parser,
		stdinR: os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		currdir: cwd,
		workdir: cwd,
		// TODO, clean this up?
		logger: parser.GetLogger(),
	}

	// RT.setupLogger()

	S, err := parser.ParseScript(args[0])
	if err != nil {
		RT.logger.Error(err)
		return err
	}

	RT.script = S
	S.Args = args[1:]

	RT.setupEnv()

	err = RT.Run()
	if err != nil {
		RT.logger.Error(err)
		return err
	}

	// cleanup

	return nil
}

func (RT *Runtime) Run() (E error) {
	R := ast.NewResult(&ast.Phase{}, nil)
	// start result
	R.BegTime = time.Now()
	defer func() {
		if R.EndTime.IsZero() {
			R.EndTime = time.Now()
		}
	}()

	for _, ph := range RT.script.Phases {
		r, err := RT.RunPhase(ph, nil)
		R.AddResult(r)
		if err != nil {
			R.AddError(err)
			break
		}

	}

	R.EndTime = time.Now()

	if len(R.Errors) == 0 {
		R.Status = 0
	} else {
		R.Status = 1
		E = fmt.Errorf("%d Script errors occurred", len(R.Errors))
	}

	// TODO print none/some/all/etc... based on config
	fmt.Fprintf(RT.Stdout, "[script]  [status:%d] [time:%v]\n", R.Status, R.EndTime.Sub(R.BegTime))

	for _, e := range R.Errors {
		RT.logger.Error(e)
	}

	return E
}
