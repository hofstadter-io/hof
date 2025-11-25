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

type ReadDirArgs struct {
	Path string `json:"path"` // path to a directory
}
type ReadDirResult struct {
	Path   string `json:"path"`   // path to a directory
	Status string `json:"status"` // "ok" or "error"
	Error  string `json:"error,omitempty"`
}

func NewReadDir() (tool.Tool, error) {
	handler := func(ctx tool.Context, input ReadDirArgs) (ReadDirResult, error) {
		state := maps.Collect(ctx.State().All())
		fmt.Println("read_dir:", input.Path, state)
		env, ok := state["env"].(map[string]any)
		if !ok {
			err := fmt.Errorf("failed to get env")
			return ReadDirResult{Path: input.Path, Status: "error", Error: err.Error()}, err
		}
		wsDir, ok := env["workspaceDir"].(string)
		if !ok {
			err := fmt.Errorf("failed to get wsDir")
			return ReadDirResult{Path: input.Path, Status: "error", Error: err.Error()}, err
		}

		entries, err := os.ReadDir(filepath.Join(wsDir, input.Path))
		if err != nil {
			return ReadDirResult{Path: input.Path, Status: "error", Error: err.Error()}, err
		}

		b := new(strings.Builder)
		for _, e := range entries {
			fmt.Fprintln(b, e.Name())
		}

		ctx.State().Set("fs:dir:"+input.Path, b.String())
		return ReadDirResult{Status: "ok", Path: input.Path}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        "read_dir",
		Description: "returns the listing of a directory",
	}, handler)
}
