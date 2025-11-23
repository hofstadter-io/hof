package filesys

import (
	"fmt"
	"io"
	"os"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/hofstadter-io/hof/lib/yagu"
)

type GlobFilesArgs struct {
	Globs []string `json:"globs"` // list of filepath globs to read into context
}
type GlobFilesResult struct {
	Files string `json:"files"`
	Error string `json:"error,omitempty"`
}

func NewGlobFiles() (tool.Tool, error) {
	handler := func(ctx tool.Context, input GlobFilesArgs) (GlobFilesResult, error) {
		fmt.Println("glob_files:", input.Globs)
		entries, err := yagu.FilepathsFromGlobs(input.Globs)
		if err != nil {
			return GlobFilesResult{Error: err.Error()}, err
		}

		b := new(strings.Builder)
		for _, entry := range entries {
			r, err := os.Open(entry)
			if err != nil {
				return GlobFilesResult{Error: err.Error()}, err
			}
			fmt.Printf("filename: %s\n", entry)
			io.Copy(b, r)
			fmt.Printf("\n\n\n")
		}

		return GlobFilesResult{Files: b.String()}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        "glob_files",
		Description: "returns the content for all files matched by a list of filepath globs",
	}, handler)
}
