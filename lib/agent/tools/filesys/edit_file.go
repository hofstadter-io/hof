package filesys

import (
	"fmt"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type EditFileArgs struct {
	Path string `json:"path"` // path to a file
	Old  string `json:"old_string"`
	New  string `json:"new_string"`
	Exp  int    `json:"expected replacements"`
}

type EditFileResult struct {
	Path   string `json:"path"`            // path to a file
	Status string `json:"status"`          // "ok" or "error"
	Error  string `json:"error,omitempty"` // an error message
}

func editFileError(path string, err error) EditFileResult {
	return EditFileResult{Path: path, Status: "error", Error: err.Error()}
}

func NewEditFile(name, prompt string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input EditFileArgs) (EditFileResult, error) {
		// need to cascade where we look, from closest to outer most
		// 1. state
		// 2. filesystem (host vs ephemeral) (tied to terminal access/env)
		// same for most tools
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Path)

		curr, err := ctx.State().Get(k)
		if err != nil {
			return editFileError(input.Path, err), err
		}

		content, ok := curr.(string)
		if !ok {
			err = fmt.Errorf("%s was not found in the cache", input.Path)
			return editFileError(input.Path, err), err
		}

		count := input.Exp
		if count == 0 {
			count = 1
		}
		found := strings.Count(content, input.Old)
		if found != count {
			err = fmt.Errorf("while editing %s, expected %d matches, but found %d", input.Path, count, found)
			return editFileError(input.Path, err), err
		}

		next := strings.Replace(content, input.Old, input.New, count)

		ctx.State().Set(k, next)
		// return the result
		return EditFileResult{Path: input.Path, Status: "ok"}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(prompt),
	}, handler)
}
