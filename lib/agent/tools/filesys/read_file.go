package filesys

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

const ReadFileTool = "read_file"
const ReadFileDesc = `
read the file at 'path' and returns the 'content'
`

type ReadFileArgs struct {
	Path string `json:"path"` // path to a file
}
type ReadFileResult struct {
	Path   string `json:"path"`            // path to a file
	Status string `json:"status"`          // "ok" or "error"
	Error  string `json:"error,omitempty"` // an error message
}

func readFileError(path string, err error) ReadFileResult {
	return ReadFileResult{Path: path, Status: "error", Error: err.Error()}
}

func NewReadFile() (tool.Tool, error) {
	handler := func(ctx tool.Context, input ReadFileArgs) (ReadFileResult, error) {
		// need to cascade where we look, from closest to outer most
		// 1. state
		// 2. artifacts
		// 3. filesystem (host vs ephemeral) (tied to terminal access/env)
		// same for most tools

		state := maps.Collect(ctx.State().All())
		fmt.Println("read_file:", input.Path, state)

		// get the env from state
		env, ok := state["env"].(map[string]any)
		if !ok {
			err := fmt.Errorf("while finding env state")
			return readFileError(input.Path, err), err
		}
		// get the workspace dir from env
		wsDir, ok := env["workspaceDir"].(string)
		if !ok {
			err := fmt.Errorf("while finding workspace directory")
			return readFileError(input.Path, err), err
		}

		// read file relative to workspace dir
		c, err := os.ReadFile(filepath.Join(wsDir, input.Path))
		if err != nil {
			fmt.Println("while reading from filesystem:", err)
			return readFileError(input.Path, err), err
		}
		ctx.State().Set("fs:file:"+input.Path, string(c))
		// return the result
		return ReadFileResult{Path: input.Path, Status: "ok"}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        ReadFileTool,
		Description: strings.TrimSpace(ReadFileDesc),
	}, handler)
}
