package environ

import (
	"fmt"
	"time"

	"dagger.io/dagger"
)

type ExecResponse struct {
	ExitCode int
	Stdout   string
	Stderr   string
	EnvID    string // resulting EnvID (not persisted)
}

func (le *localEnviron) Exec(envUri, nextTag, script string) (resp ExecResponse, err error) {

	_, env, err := le.lookupEnviron(envUri)
	if err != nil {
		return resp, fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	path, err := extractPathEmptyOk(envUri)
	if err != nil {
		return resp, fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	runner := env.
		WithEnvVariable("CGO_ENABLED", "1")
		// TODO, this needs to be defined on the outside, project specific
		// TODO, this should be defined on the env itself, maybe we need more in the moment? (can just do itself for now)
		// inject user ENV & SHH vars

	// pushd/popd - part 1 - "pushd"
	cwd, err := runner.Workdir(le.ctx)
	if err != nil {
		return resp, fmt.Errorf("while getting working directory(%s): %w", envUri, err)
	}
	if path != "" {
		runner = runner.WithWorkdir(path)
	}

	scriptHeader := `
#!/bin/bash
set -euo pipefail

`
	fullScript := scriptHeader + script
	fmt.Printf("Running script in %s:%s:%s\n", envUri, path)
	fmt.Println(fullScript)

	//
	// This is where we actually run the exec
	//
	//    TODO, sequences of exec should persist
	result, err := runner.WithEnvVariable("CACHE_BUST", time.Now().Local().String()).
		WithExec([]string{"sh", "-c", script}, dagger.ContainerWithExecOpts{
			Expect:         dagger.ReturnTypeAny,
			RedirectStdout: "/stdout.txt",
			RedirectStderr: "/stderr.txt",
			Expand:         true,
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

	// pushd/popd - part 2 - "popd"
	if path != "" {
		result = result.WithWorkdir(cwd)
	}

	// get final ID
	envId, err := result.ID(le.ctx)
	if err != nil {
		return resp, fmt.Errorf("while getting envId: %w", err)
	}
	resp.EnvID = string(envId)

	// TODO this should move out
	// 1. we don't always want to save it (ephemeral)
	// 2. the caller should decide what/when
	//
	// persist exec
	//
	nextUri := replaceTag(envUri, nextTag)
	err = le.persistEnviron(nextUri, nil, runner)

	return resp, nil
}
