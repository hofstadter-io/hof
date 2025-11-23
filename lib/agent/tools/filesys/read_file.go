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
	Status       string `json:"status"`                  // "ok" or "error"
	Path         string `json:"path"`                    // path to a file
	Content      string `json:"content,omitempty"`       // content of the file
	ErrorMessage string `json:"error_message,omitempty"` // an error message
}

func readFileError(err error) ReadFileResult {
	return ReadFileResult{Status: "error", ErrorMessage: err.Error()}
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
			return readFileError(err), err
		}
		// get the workspace dir from env
		wsDir, ok := env["workspaceDir"].(string)
		if !ok {
			err := fmt.Errorf("while finding workspace directory")
			return readFileError(err), err
		}

		// read file relative to workspace dir
		c, err := os.ReadFile(filepath.Join(wsDir, input.Path))
		if err != nil {
			fmt.Println("while reading from filesystem:", err)
			return readFileError(err), err
		}
		// return the result
		return ReadFileResult{Content: string(c), Status: "ok"}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        ReadFileTool,
		Description: strings.TrimSpace(ReadFileDesc),
	}, handler)
}
