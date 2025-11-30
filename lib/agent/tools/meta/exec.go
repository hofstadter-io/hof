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
	Key    string `json:"key"`    // key for the cache entry
	Script string `json:"script"` // command or script to run
}
type ExecResult struct {
	Key    string `json:"key"`             // path to a cache entry
	Status string `json:"status"`          // "ok" or "error"
	Error  string `json:"error,omitempty"` // error message if there is an error
}

func execError(key string, err error) ExecResult {
	fmt.Println("ERROR:", key, err)
	return ExecResult{Key: key, Status: "error", Error: err.Error()}
}

func Exec(name, description, runenv string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input ExecArgs) (ExecResult, error) {
		// calculate our real key
		k := fmt.Sprintf("%s:%s", ctx.AgentName(), input.Key)
		fmt.Printf("%s:%s\n", name, k)

		// workdir is always set by us
		w, _ := ctx.State().Get("basedir")
		workdir := w.(string)

		//
		// Get the latest Dagger layer (dir)
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
			WithEnvVariable("CGO_ENABLED", "1").
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
			return execError(input.Key, err), nil
		}

		//
		// collect results
		//
		exitCode, err := result.ExitCode(ctx)
		if err != nil {
			return execError(input.Key, err), nil
		}
		stdout, err := result.File("/stdout.txt").Contents(ctx)
		if err != nil {
			return execError(input.Key, err), nil
		}
		stderr, err := result.File("/stderr.txt").Contents(ctx)
		if err != nil {
			return execError(input.Key, err), nil
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
			return execError(input.Key, err), nil
		}
		// newId, err := result.Directory(workdir).ID(ctx)
		// if err != nil {
		// 	return execError(input.Key, err), nil
		// }

		//
		// update state
		//
		outputFmt := "ExitCode: %d\n\nStdout:\n%s\n\nStderr:\n%s\n\n"
		err = ctx.State().Set(k, fmt.Sprintf(outputFmt, exitCode, stdout, stderr))
		if err != nil {
			return execError(input.Key, err), nil
		}
		// err = ctx.State().Set("dagger", string(newId))
		// if err != nil {
		// 	return execError(input.Key, err), nil
		// }
		err = ctx.State().Set("runenv", string(envId))
		if err != nil {
			return execError(input.Key, err), nil
		}

		if exitCode != 0 {
			return execError(input.Key, fmt.Errorf("command failed with exit code: %d", exitCode)), nil
		}

		// return status result
		return ExecResult{Status: "ok", Key: input.Key}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}
