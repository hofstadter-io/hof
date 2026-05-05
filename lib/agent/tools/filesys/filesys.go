package filesys

import (
	"fmt"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/hofstadter-io/hof/lib/agent/services/environ"
)

type FilesysPathArgs struct {
	Path string `json:"path"` // filesystem path
}

type FilesysResult struct {
	Path   string `json:"path"`            // filesystem path
	Status string `json:"status"`          // "ok" or "error"
	Error  string `json:"error,omitempty"` // the error message if present
}

type FilesysOutputResult struct {
	Path   string `json:"path"`            // filesystem path
	Output string `json:"output"`          // output of the filesystem query
	Status string `json:"status"`          // "ok" or "error"
	Error  string `json:"error,omitempty"` // the error message if present
}

func filesysError(path string, err error) FilesysResult {
	fmt.Println("Filesys.ERROR:", path, err)
	return FilesysResult{Path: path, Status: "error", Error: err.Error()}
}
func filesysOutputError(path string, err error) FilesysOutputResult {
	fmt.Println("Filesys.ERROR:", path, err)
	return FilesysOutputResult{Path: path, Status: "error", Error: err.Error()}
}

func getAndCheckCurrEnv(ctx tool.Context) (string, error) {
	// Get environ
	currUri, err := ctx.State().Get("currEnv")
	if err != nil {
		return "", fmt.Errorf("while getting currEnv from state: %w", err)
	}
	if currUri == nil {
		return "", fmt.Errorf("no environment, attach a filesystem or container")
	}

	currStr, ok := currUri.(string)
	if !ok {
		return "", fmt.Errorf("state.currEnv is not a string")
	}

	return currStr, nil
}

func FilesysRead(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysPathArgs) (FilesysResult, error) {
		// workdir is always set by us
		// w, _ := ctx.State().Get("basedir")
		// workdir := w.(string)
		// TODO, perhaps some cleaning or checking it is not an absolute path while constucting the real path

		// calculate our real key
		k := fmt.Sprintf("files:%s:%s", ctx.AgentName(), input.Path)
		fmt.Printf("%s:%s\n", name, k)

		// Get environ
		currUri, err := getAndCheckCurrEnv(ctx)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// TODO path shenanigans
		fmt.Printf("%s:%s @ %s\n", name, k, currUri)

		// Read file content
		content, err := environ.Client().ReadFile(currUri, input.Path, false)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// Add to State
		err = ctx.State().Set(k, content)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// return status result
		return FilesysResult{Path: input.Path, Status: "ok"}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

func FilesysList(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysPathArgs) (FilesysOutputResult, error) {
		// workdir is always set by us
		// w, _ := ctx.State().Get("basedir")
		// workdir := w.(string)
		// TODO, perhaps some cleaning or checking it is not an absolute path while constucting the real path

		// Get environ
		currUri, err := getAndCheckCurrEnv(ctx)
		if err != nil {
			return filesysOutputError(input.Path, err), nil
		}

		// TODO path shenanigans
		fmt.Printf("%s:%s @ %s\n", name, input.Path, currUri)

		// Get directory list
		dirList, err := environ.Client().ReadDirectory(currUri, input.Path, false)
		if err != nil {
			return filesysOutputError(input.Path, err), nil
		}

		// Construct output string
		b := new(strings.Builder)
		for _, e := range dirList.Entries {
			fmt.Fprintln(b, e.Name)
		}

		// return status result
		return FilesysOutputResult{Status: "ok", Path: input.Path, Output: b.String()}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

type FilesysGlobArgs struct {
	Glob string `json:"glob"` // pattern to match (e.g. "*.md")
}

func FilesysGlob(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysGlobArgs) (FilesysOutputResult, error) {
		// workdir is always set by us
		// w, _ := ctx.State().Get("basedir")
		// workdir := w.(string)
		// TODO, perhaps some cleaning or checking it is not an absolute path while constucting the real path

		// Get environ
		currUri, err := getAndCheckCurrEnv(ctx)
		if err != nil {
			return filesysOutputError(input.Glob, err), nil
		}

		// TODO path shenanigans ?

		// grep directory
		results, err := environ.Client().GlobDirectory(currUri, input.Glob, false)
		if err != nil {
			return filesysOutputError(input.Glob, err), nil
		}

		if len(results) > 100 {
			return filesysOutputError(input.Glob, fmt.Errorf("too many results, narrow your glob pattern, suggest to limit depth by not using double star or limit to a subdirectory by using it as the prefix")), nil
		}

		// build output message
		b := new(strings.Builder)
		for _, r := range results {
			fmt.Fprintf(b, "%s\n", r)
		}

		// return the result
		return FilesysOutputResult{Status: "ok", Path: input.Glob, Output: b.String()}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

type FilesysGrepArgs struct {
	Path   string `json:"path"`   // base path to grep from
	Glob   string `json:"glob"`   // pattern to match (e.g. "*.md")
	Regexp string `json:"regexp"` // a regular expression to grep for
}

func FilesysGrep(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysGrepArgs) (FilesysOutputResult, error) {
		// workdir is always set by us
		// w, _ := ctx.State().Get("basedir")
		// workdir := w.(string)
		// TODO, perhaps some cleaning or checking it is not an absolute path while constucting the real path

		// Get environ
		currUri, err := getAndCheckCurrEnv(ctx)
		if err != nil {
			return filesysOutputError(input.Path, err), nil
		}

		// TODO path shenanigans

		// grep directory
		results, err := environ.Client().GrepDirectory(currUri, input.Regexp, false)
		if err != nil {
			return filesysOutputError(input.Path, err), nil
		}

		// build output message
		b := new(strings.Builder)
		for _, r := range results {
			fp, _ := r.FilePath(ctx)
			ln, _ := r.LineNumber(ctx)
			ml, _ := r.MatchedLines(ctx)
			fmt.Fprintf(b, "%s:%d:%s\n", fp, ln, ml)
		}

		// return the result
		return FilesysOutputResult{Status: "ok", Path: input.Path, Output: b.String()}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

type FilesysEditArgs struct {
	Path  string           `json:"path"` // path to a file
	Edits []environ.EditOp `json:"edits"`
}

func FilesysEdit(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysEditArgs) (FilesysResult, error) {
		k := fmt.Sprintf("files:%s:%s", ctx.AgentName(), input.Path)
		fmt.Printf("fsEdit.key: %s:%s\n", name, k)

		// get the current env, it's in Uri format
		currUri, err := getAndCheckCurrEnv(ctx)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// edit the file
		nextUri, nextContent, err := environ.Client().EditFile(currUri, input.Path, input.Edits)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		//
		// write to ADK state
		//
		err = ctx.State().Set(k, nextContent)
		if err != nil {
			return filesysError(input.Path, err), nil
		}
		err = ctx.State().Set("currEnv", nextUri)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// return status result
		return FilesysResult{Status: "ok", Path: input.Path}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

type FilesysWriteArgs struct {
	Path    string `json:"path"`    // path to a directory
	Content string `json:"content"` // path to a directory
}

func FilesysWrite(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysWriteArgs) (FilesysResult, error) {
		// calculate our real key
		k := fmt.Sprintf("files:%s:%s", ctx.AgentName(), input.Path)
		fmt.Printf("%s:%s\n", name, k)

		// get the current env, it's in Uri format
		currUri, err := getAndCheckCurrEnv(ctx)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// write the file
		nextUri, err := environ.Client().WriteFile(currUri, input.Path, input.Content)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		//
		// write to ADK state
		//
		err = ctx.State().Set(k, input.Content)
		if err != nil {
			return filesysError(input.Path, err), nil
		}
		err = ctx.State().Set("currEnv", nextUri)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// return status result
		return FilesysResult{Status: "ok", Path: input.Path}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

func FilesysMkdir(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysPathArgs) (FilesysResult, error) {
		// get the current env, it's in Uri format
		currUri, err := getAndCheckCurrEnv(ctx)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// create the directory
		nextUri, err := environ.Client().CreateDirectory(currUri, input.Path)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// update state
		err = ctx.State().Set("currEnv", nextUri)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// return status result
		return FilesysResult{Status: "ok", Path: input.Path}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

type FilesysMoveArgs struct {
	Src string `json:"src"` // source path
	Dst string `json:"dst"` // destination path
}

func FilesysRename(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysMoveArgs) (FilesysResult, error) {
		// get the current env, it's in Uri format
		currUri, err := getAndCheckCurrEnv(ctx)
		if err != nil {
			return filesysError(input.Src, err), nil
		}

		// move the path
		nextUri, err := environ.Client().Move(currUri, input.Src, input.Dst, true)
		if err != nil {
			return filesysError(input.Src, err), nil
		}

		// update state
		err = ctx.State().Set("currEnv", nextUri)
		if err != nil {
			return filesysError(input.Src, err), nil
		}

		// return status result
		return FilesysResult{Status: "ok", Path: input.Src}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

func FilesysCopy(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysMoveArgs) (FilesysResult, error) {
		// get the current env, it's in Uri format
		currUri, err := getAndCheckCurrEnv(ctx)
		if err != nil {
			return filesysError(input.Src, err), nil
		}

		// copy the path
		nextUri, err := environ.Client().Copy(currUri, input.Src, input.Dst, true)
		if err != nil {
			return filesysError(input.Src, err), nil
		}

		// update state
		err = ctx.State().Set("currEnv", nextUri)
		if err != nil {
			return filesysError(input.Src, err), nil
		}

		// return status result
		return FilesysResult{Status: "ok", Path: input.Src}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

func FilesysDel(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysPathArgs) (FilesysResult, error) {
		// calculate our real key
		f := fmt.Sprintf("files:%s:%s", ctx.AgentName(), input.Path)
		fmt.Printf("%s:%s\n", name, f)
		k := fmt.Sprintf("files:%s:%s", ctx.AgentName(), input.Path)
		fmt.Printf("%s:%s\n", name, k)

		// get the current env, it's in Uri format
		currUri, err := getAndCheckCurrEnv(ctx)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// write the file
		nextUri, err := environ.Client().Delete(currUri, input.Path, true)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		//
		// write to ADK state
		//
		err = ctx.State().Set(f, nil)
		if err != nil {
			return filesysError(input.Path, err), nil
		}
		err = ctx.State().Set(k, nil)
		if err != nil {
			return filesysError(input.Path, err), nil
		}
		err = ctx.State().Set("currEnv", nextUri)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// return status result
		return FilesysResult{Status: "ok", Path: input.Path}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}
