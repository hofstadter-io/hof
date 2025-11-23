package filesys

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type PatchFilesArgs struct {
	Patches string `json:"patches"` // a series of git patches to apply
}
type PatchFilesResult struct {
	Paths  []string `json:"paths"`
	Status string   `json:"status"`
	Error  string   `json:"error,omitempty"`
}

func NewPatchFiles() (tool.Tool, error) {
	handler := func(ctx tool.Context, input PatchFilesArgs) (PatchFilesResult, error) {
		fmt.Println("patch_file:", input.Patches)
		var r PatchFilesResult

		files, _, err := gitdiff.Parse(strings.NewReader(input.Patches))
		if err != nil {
			r.Error = err.Error()
			return r, nil
		}

		for _, file := range files {
			f, err := os.OpenFile(file.OldName, os.O_RDONLY, 0o644)
			if err != nil {
				r.Error = err.Error()
				return r, nil
			}

			var output bytes.Buffer
			err = gitdiff.Apply(&output, f, files[0])
			if err != nil {
				r.Error = err.Error()
				return r, nil
			}
			err = os.WriteFile(file.OldName, output.Bytes(), 0o644)
			if err != nil {
				r.Error = err.Error()
				return r, nil
			}
			r.Paths = append(r.Paths, file.OldName)
		}

		r.Status = "ok"
		return r, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        "patch_files",
		Description: "apply a set of patches to one or more files using git or diff format",
	}, handler)
}
