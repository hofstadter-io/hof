package exec

import (
	"fmt"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/hofstadter-io/hof/lib/agent/services/environ"
	"github.com/kr/pretty"
)

type ExecArgs struct {
	Script string `json:"script"` // command or script to run
}
type ExecResult struct {
	ExitCode int    `json:"exitCode"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Status   string `json:"status"`          // "ok" or "error"
	Error    string `json:"error,omitempty"` // error message if there is an error
}

func getAndCheckCurrEnv(ctx tool.Context) (string, error) {
	// Get environ
	currUri, err := ctx.State().Get("currEnv")
	if err != nil {
		return "", fmt.Errorf("while getting currEnv from state: %w", err)
	}
	if currUri == nil {
		return "", fmt.Errorf("no environment, attach a filesystem or container")
	}

	currStr, ok := currUri.(string)
	if !ok {
		return "", fmt.Errorf("state.currEnv is not a string")
	}

	return currStr, nil
}

func execError(err error) ExecResult {
	fmt.Println("EXEC.error:", err)
	return ExecResult{Status: "error", Error: err.Error()}
}

func Exec(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input ExecArgs) (ExecResult, error) {
		// calculate our real key
		k := fmt.Sprintf("%s:%s", ctx.AgentName(), input.Script[:min(42, len(input.Script))])
		fmt.Printf("%s:%s\n", name, k)

		// get the current env, it's in Uri format
		currUri, err := getAndCheckCurrEnv(ctx)
		if err != nil {
			return execError(err), nil
		}

		resp, err := environ.Client().Exec(currUri, input.Script)
		if err != nil {
			return execError(err), nil
		}

		status := "ok"
		if resp.ExitCode != 0 {
			status = "error"
		}

		// update state
		err = ctx.State().Set("currEnv", resp.NextUri)
		if err != nil {
			return execError(err), nil
		}

		// return exec results
		final := ExecResult{Status: status, ExitCode: resp.ExitCode, Stdout: resp.Stdout, Stderr: resp.Stderr}
		fmt.Printf("final: %#+v\n", pretty.Formatter(final))
		return final, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}
