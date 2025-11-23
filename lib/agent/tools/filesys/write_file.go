package filesys

import (
	"errors"
	"fmt"

	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type WriteFileArgs struct {
	Path    string `json:"path"`    // path to a file
	Content string `json:"content"` // contents to write to file
}
type WriteFileResult struct {
	Path   string `json:"path"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func NewWriteFile() (tool.Tool, error) {
	handler := func(ctx tool.Context, input WriteFileArgs) (WriteFileResult, error) {
		fmt.Println("write_file:", input.Path)

		// we just add the contents to the state and handle on the frontend
		// in the end, we need to mirror the current "accepted" state
		// in both the frontend and this backend (when context constructing)
		val, err := ctx.State().Get("fs")
		if err != nil && !errors.Is(err, session.ErrStateKeyNotExist) {
			return WriteFileResult{Path: input.Path, Status: "error", Error: err.Error()}, nil
		}
		var fs map[string]string
		if fs == nil {
			fs = make(map[string]string)
		} else {
			fs = val.(map[string]string)
		}

		// we only keep the most recent version in state, we can update based on the user's final choices / checkpoints
		fs[input.Path] = input.Content
		ctx.State().Set("fs", fs)

		return WriteFileResult{Path: input.Path, Status: "ok"}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        "write_file",
		Description: "Suggests changes to a file and present the user with a diff. The user will accept or reject, whole or parts, and then save the file.",
	}, handler)
}
