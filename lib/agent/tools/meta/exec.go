// this package is really the next iteration, filesys is legacy

package meta

import (
	"fmt"
	"strings"
	"time"

	"dagger.io/dagger"
	vegdagger "github.com/hofstadter-io/hof/lib/agent/runtime/dagger"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type ExecArgs struct {
	Script string `json:"script"` // command or script to run
}
type ExecResult struct {
	ExitCode int    `json:"exitCode,omitempty"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Status   string `json:"status"`          // "ok" or "error"
	Error    string `json:"error,omitempty"` // error message if there is an error
}

func execError(err error) ExecResult {
	fmt.Println("EXEC.error:", err)
	return ExecResult{Status: "error", Error: err.Error()}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Exec(name, description, runenv string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input ExecArgs) (ExecResult, error) {
		// calculate our real key
		k := fmt.Sprintf("%s:%s", ctx.AgentName(), input.Script[:min(42, len(input.Script))])
		fmt.Printf("%s:%s\n", name, k)

		// workdir is always set by us
		w, _ := ctx.State().Get("basedir")
		workdir := w.(string)

		//
		// Get the latest Dagger layer (dir or container?)
		//
		// get the client
		dag, _ := vegdagger.Get(ctx)
		// get the directory
		dagId, _ := ctx.State().Get("dagger")
		dir := dag.LoadDirectoryFromID(dagger.DirectoryID(dagId.(string)))
		wdir := dir.Directory(workdir)
		// get the runenv
		// envId, _ := ctx.State().Get("runenv") // this might be a directory or container now...
		// cnt := dag.LoadContainerFromID(dagger.ContainerID(envId.(string)))

		// hmm, attaching vs from start
		// directory vs container lookup
		// we also don't know the path, mounting might be better tbh, but where
		// mounting makes writing it back more challenging, or we lose the optimization?
		// need to do some path finangling, and make sure the agent doesn't screw it up, so keep simple for now
		container := dag.Container().From(runenv).
			WithEnvVariable("CGO_ENABLED", "1"). // TODO, this needs to be defined on the outside, project specific
			WithWorkdir(workdir)
		runner := container.WithMountedDirectory(workdir, wdir)

		scriptHeader := `
#!/bin/bash
set -euo pipefail

`
		script := scriptHeader + input.Script
		fmt.Println("Running script in :", workdir)
		fmt.Println(script)

		//
		// This is where we actually run the exec
		//
		//    TODO, sequences of exec should persist
		result, err := runner.WithEnvVariable("CACHE_BUST", time.Now().Local().String()).
			WithExec([]string{"sh", "-c", script}, dagger.ContainerWithExecOpts{
				Expect:         dagger.ReturnTypeAny,
				RedirectStdout: "/stdout.txt",
				RedirectStderr: "/stderr.txt",
			}).Sync(ctx)
		if err != nil {
			fmt.Println("DONT WANT THIS ERROR:", err)
			return execError(err), nil
		}

		//
		// collect results
		//
		exitCode, err := result.ExitCode(ctx)
		if err != nil {
			return execError(err), nil
		}
		stdout, err := result.File("/stdout.txt").Contents(ctx)
		if err != nil {
			return execError(err), nil
		}
		stderr, err := result.File("/stderr.txt").Contents(ctx)
		if err != nil {
			return execError(err), nil
		}

		fmt.Println("ExitCode:", exitCode)
		fmt.Println("Stdout:")
		fmt.Println(stdout)
		fmt.Println("Stderr:")
		fmt.Println(stderr)

		//
		// new ids
		//
		envId, err := result.ID(ctx)
		if err != nil {
			return execError(err), nil
		}
		// newId, err := result.Directory(workdir).ID(ctx)
		// if err != nil {
		// 	return execError(input.Key, err), nil
		// }

		//
		// update state (filesys,execenv)
		//
		// err = ctx.State().Set("dagger", string(newId))
		// if err != nil {
		// 	return execError(input.Key, err), nil
		// }
		err = ctx.State().Set("runenv", string(envId))
		if err != nil {
			return execError(err), nil
		}

		status := "ok"
		if exitCode != 0 {
			status = "error"
		}

		// return status result
		return ExecResult{Status: status, ExitCode: exitCode, Stdout: stdout, Stderr: stderr}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}
