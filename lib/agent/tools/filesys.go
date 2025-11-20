package tools

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type ReadFileArgs struct {
	Path string `json:"path"` // path to a file
}
type ReadFileResult struct {
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

func NewReadFile() (tool.Tool, error) {
	handler := func(ctx tool.Context, input ReadFileArgs) (ReadFileResult, error) {
		fmt.Println("read_file:", input.Path)
		c, err := os.ReadFile(input.Path)
		if err != nil {
			return ReadFileResult{Error: err.Error()}, err
		}
		return ReadFileResult{Content: string(c)}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        "read_file",
		Description: "returns the content of a file",
	}, handler)
}

type ReadDirArgs struct {
	Path string `json:"path"` // path to a directory
}
type ReadDirResult struct {
	Listing string `json:"listing"`
	Error   string `json:"error,omitempty"`
}

func NewReadDir() (tool.Tool, error) {
	handler := func(ctx tool.Context, input ReadDirArgs) (ReadDirResult, error) {
		fmt.Println("read_dir:", input.Path)
		entries, err := os.ReadDir(input.Path)
		if err != nil {
			return ReadDirResult{Error: err.Error()}, err
		}

		b := new(strings.Builder)
		for _, e := range entries {
			fmt.Fprintln(b, e.Name())
		}

		return ReadDirResult{Listing: b.String()}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        "read_dir",
		Description: "returns the listing of a directory",
	}, handler)
}

type TreeDirArgs struct {
	Path         string `json:"path"`         // path to a file
	IncludeFiles bool   `json:"includeFiles"` // include files in the tree listing
}
type TreeDirResult struct {
	Listing string `json:"listing"`
	Error   string `json:"error"`
}

func NewTreeDir() (tool.Tool, error) {
	handler := func(ctx tool.Context, input TreeDirArgs) (TreeDirResult, error) {
		fmt.Println("tree_dir:", input.Path, input.IncludeFiles)
		b := new(strings.Builder)

		err := walkDir(input.Path, input.IncludeFiles, b, "")
		if err != nil {
			return TreeDirResult{Error: err.Error()}, err
		}

		return TreeDirResult{Listing: b.String()}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        "tree_dir",
		Description: "returns the tree layout of a directory, leading hyphens determine depth",
	}, handler)
}

func walkDir(path string, includeFiles bool, b *strings.Builder, prefix string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			fmt.Fprintf(b, "%s%s", prefix, e.Name())
			err = walkDir(filepath.Join(path, e.Name()), includeFiles, b, prefix+"-")
			if err != nil {
				return err
			}
		}
		if includeFiles {
			fmt.Fprintf(b, "%s%s", prefix, e.Name())
		}
	}

	return nil
}

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
		dir := filepath.Dir(input.Path)
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			return WriteFileResult{Path: input.Path, Error: err.Error()}, err
		}
		err = os.WriteFile(input.Path, []byte(input.Content), 0o644)
		if err != nil {
			return WriteFileResult{Path: input.Path, Error: err.Error()}, err
		}

		return WriteFileResult{Path: input.Path, Status: "ok"}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        "write_file",
		Description: "writes content to a file, replacing any existing content",
	}, handler)
}

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
