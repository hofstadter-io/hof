package environ

import (
	"fmt"

	"dagger.io/dagger"
)

type ExecRequest struct {
	Script  string
	Workdir string
	// ... what else might we want to expose generally and to an agent
}

type ExecResponse struct {
	ExitCode int
	Stdout   string
	Stderr   string
	NextUri  string
}

func (le *localEnviron) Exec(envUri, script string) (resp ExecResponse, err error) {

	_, env, err := le.LookupEnviron(envUri)
	if err != nil {
		return resp, fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	// path, err := extractPathEmptyOk(envUri)
	// if err != nil {
	// 	return resp, fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	// }

	runner := env.
		WithEnvVariable("CGO_ENABLED", "1")
		// TODO, this needs to be defined on the outside, project specific
		// TODO, this should be defined on the env itself, maybe we need more in the moment? (can just do itself for now)
		// inject user ENV & SHH vars

	// if workdir not set
	cwd, err := runner.Workdir(le.ctx)
	if err != nil {
		return resp, fmt.Errorf("while getting working directory(%s): %w", envUri, err)
	}

	scriptHeader := `
#!/bin/sh
set -euo pipefail

`
	fullScript := scriptHeader + script
	fmt.Printf("Running script in %s at %s\n", envUri, cwd)
	fmt.Println(fullScript)

	//
	// This is where we actually run the exec
	//
	//    TODO, sequences of exec should persist
	result, err := runner.
		// WithEnvVariable("CACHE_BUST", time.Now().Local().String()).
		WithExec([]string{"sh", "-c", fullScript}, dagger.ContainerWithExecOpts{
			Expect:         dagger.ReturnTypeAny,
			RedirectStdout: "/stdout.txt",
			RedirectStderr: "/stderr.txt",
			Expand:         true,
			// Workdir:        path,
		}).Sync(le.ctx)
	if err != nil {
		return resp, fmt.Errorf("while running script: %w", err)
	}

	//
	// collect results
	//
	resp.ExitCode, err = result.ExitCode(le.ctx)
	if err != nil {
		return resp, fmt.Errorf("while getting exit code: %w", err)
	}
	resp.Stdout, err = result.File("/stdout.txt").Contents(le.ctx)
	if err != nil {
		return resp, fmt.Errorf("while getting stdout: %w", err)
	}
	resp.Stderr, err = result.File("/stderr.txt").Contents(le.ctx)
	if err != nil {
		return resp, fmt.Errorf("while getting stderr: %w", err)
	}

	nextUri, _, err := IncrementTag(envUri)
	if err != nil {
		return resp, fmt.Errorf("while incrementing tag(%s): %w", envUri, err)
	}

	fmt.Println("Environ.persisting", nextUri, result)

	err = le.persistEnviron(nextUri, nil, result)
	if err != nil {
		return resp, fmt.Errorf("while persisting environ(%s -> %s): %w", envUri, nextUri, err)
	}
	resp.NextUri = nextUri

	return resp, nil
}
