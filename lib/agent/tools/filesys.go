package tools

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/hofstadter-io/hof/lib/yagu"
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
