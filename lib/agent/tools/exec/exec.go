package exec

import (
	"fmt"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/hofstadter-io/hof/lib/agent/runtime/services/environ"
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

func Exec(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input ExecArgs) (ExecResult, error) {
		// calculate our real key
		k := fmt.Sprintf("%s:%s", ctx.AgentName(), input.Script[:min(42, len(input.Script))])
		fmt.Printf("%s:%s\n", name, k)

		envUri, _ := ctx.State().Get("currEnv")

		resp, err := environ.Client().Exec(envUri.(string), "...", input.Script)

		if err != nil {
			return execError(err), nil
		}

		status := "ok"
		if resp.ExitCode != 0 {
			status = "error"
		}

		// TODO, persist envId
		// need to extract and update the tag, which requires figuring out all the path, tag, qp BS...

		// return status result
		return ExecResult{Status: status, ExitCode: resp.ExitCode, Stdout: resp.Stdout, Stderr: resp.Stderr}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}
