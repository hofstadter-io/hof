package filesys

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/hofstadter-io/hof/lib/yagu"
)

type TreeDirArgs struct {
	Path         string `json:"path"`            // path to a file
	IncludeFiles bool   `json:"includeFiles"`    // include files in the tree listing
	Depth        int    `json:"depth,omitempty"` // how deep to traverse, limited to 10
}
type TreeDirResult struct {
	Listing string `json:"listing"`
	Error   string `json:"error"`
}

func NewTreeDir() (tool.Tool, error) {
	handler := func(ctx tool.Context, input TreeDirArgs) (TreeDirResult, error) {
		fmt.Println("tree_dir:", input.Path, input.IncludeFiles)
		b := new(strings.Builder)

		if input.Depth < 1 {
			input.Depth = 4
		}
		if input.Depth > 10 {
			input.Depth = 10
		}

		// respect .gitignore
		// ideally only need to load this once at startup (lazily)
		ign, err := yagu.IgnFromGit()
		if err != nil {
			return TreeDirResult{Error: err.Error()}, err
		}

		// perhaps use ign walk here, and also add file including or not
		err = walkDir(input.Path, input.IncludeFiles, input.Depth, ign, b, "")
		if err != nil {
			return TreeDirResult{Error: err.Error()}, err
		}

		return TreeDirResult{Listing: b.String()}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        "tree_dir",
		Description: "Returns the tree layout of a directory, including files should be done rarely. Depth controls how deep the listing can go, defaults to 4. Tabs are used for indentation to indicate depth.",
	}, handler)
}

func walkDir(path string, includeFiles bool, depth int, ign yagu.IgnoreList, b *strings.Builder, prefix string) error {

	if ign.Match(filepath.Join(path, path)) {
		return nil
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if ign.Match(filepath.Join(path, e.Name())) {
			continue
		}
		if e.IsDir() {
			fmt.Fprintf(b, "%s%s\n", prefix, e.Name())
			err = walkDir(filepath.Join(path, e.Name()), includeFiles, depth-1, ign, b, "\t"+prefix)
			if err != nil {
				return err
			}
		}

		// TEMP not including files because node_modules keeps getting picked up
		// if includeFiles {
		// 	if ign.Match(filepath.Join(path, e.Name())) {
		// 		continue
		// 	}
		// 	fmt.Fprintf(b, "%s%s\n", prefix, e.Name())
		// }
	}

	return nil
}
