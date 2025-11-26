package meta

import (
	"fmt"
	"strings"

	"dagger.io/dagger"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	vegdagger "github.com/hofstadter-io/hof/lib/agent/runtime/dagger"
)

type FilesysPathArgs struct {
	Path string `json:"path"` // filesystem path
}

type FilesysResult struct {
	Path   string `json:"path"`            // filesystem path
	Status string `json:"status"`          // "ok" or "error"
	Error  string `json:"error,omitempty"` // the error message if present
}

func filesysError(path string, err error) FilesysResult {
	fmt.Println("Filesys.ERROR:", path, err)
	return FilesysResult{Path: path, Status: "error", Error: err.Error()}
}

func FilesysRead(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysPathArgs) (FilesysResult, error) {
		// calculate our real key
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Path)
		fmt.Printf("%s:%s\n", name, k)

		//
		// Get from Dagger
		//
		dagId, _ := ctx.State().Get("dagger")
		dag, _ := vegdagger.Get(ctx)
		dir := dag.LoadDirectoryFromID(dagger.DirectoryID(dagId.(string)))
		file, err := dir.File(input.Path).Contents(ctx)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// Add to State
		err = ctx.State().Set(k, file)
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
	handler := func(ctx tool.Context, input FilesysPathArgs) (FilesysResult, error) {
		// calculate our real key
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Path)
		fmt.Printf("%s:%s\n", name, k)

		//
		// Get from Dagger
		//
		dagId, _ := ctx.State().Get("dagger")
		dag, _ := vegdagger.Get(ctx)
		dir := dag.LoadDirectoryFromID(dagger.DirectoryID(dagId.(string)))
		entries, err := dir.Directory(input.Path).Entries(ctx)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		//
		// Add to State
		//
		b := new(strings.Builder)
		for _, e := range entries {
			fmt.Fprintln(b, e)
		}
		err = ctx.State().Set(k, b.String())
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

type FilesysGrepArgs struct {
	Path   string `json:"path"`   // base path to grep from
	Glob   string `json:"glob"`   // pattern to match (e.g. "*.md")
	Regexp string `json:"regexp"` // a regular expression to grep for
}

func FilesysGrep(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysGrepArgs) (FilesysResult, error) {
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Path)
		fmt.Printf("%s:%s\n", name, k)

		//
		// Get from Dagger
		//
		dagId, _ := ctx.State().Get("dagger")
		dag, _ := vegdagger.Get(ctx)
		dir := dag.LoadDirectoryFromID(dagger.DirectoryID(dagId.(string)))

		results, err := dir.Search(ctx, input.Regexp, dagger.DirectorySearchOpts{
			Limit:       100,
			SkipIgnored: true,
		})
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		b := new(strings.Builder)
		for _, r := range results {
			fp, _ := r.FilePath(ctx)
			ln, _ := r.LineNumber(ctx)
			ml, _ := r.MatchedLines(ctx)
			fmt.Fprintf(b, "%s:%d:%s\n", fp, ln, ml)
		}

		// write to state
		err = ctx.State().Set(k, b.String())
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// return the result
		return FilesysResult{Path: input.Path, Status: "ok"}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

type FilesysEditArgs struct {
	Path string `json:"path"` // path to a file
	Old  string `json:"old_string"`
	New  string `json:"new_string"`
	Exp  int    `json:"expected_replacements"` // defaults to 1 if not set
}

func FilesysEdit(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysEditArgs) (FilesysResult, error) {
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Path)
		fmt.Printf("%s:%s\n", name, k)

		//
		// Get from Dagger
		//
		dagId, _ := ctx.State().Get("dagger")
		dag, _ := vegdagger.Get(ctx)
		dir := dag.LoadDirectoryFromID(dagger.DirectoryID(dagId.(string)))
		content, err := dir.File(input.Path).Contents(ctx)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		//
		// Check and Replace content
		//
		count := input.Exp
		if count == 0 {
			count = 1
		}
		found := strings.Count(content, input.Old)
		if found != count {
			err = fmt.Errorf("while editing %q, expected %d matches, but found %d", input.Path, count, found)
			return filesysError(input.Path, err), nil
		}
		next := strings.Replace(content, input.Old, input.New, count)

		//
		// Update in Dagger
		//
		dir = dir.WithNewFile(input.Path, next, dagger.DirectoryWithNewFileOpts{})
		newId, err := dir.ID(ctx)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		// Add to State
		err = ctx.State().Set(k, next)
		if err != nil {
			return filesysError(input.Path, err), nil
		}
		err = ctx.State().Set("dagger", string(newId))
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

type FilesysWriteArgs struct {
	Path    string `json:"path"`    // path to a directory
	Content string `json:"content"` // path to a directory
}

func FilesysWrite(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysWriteArgs) (FilesysResult, error) {
		// calculate our real key
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Path)
		fmt.Printf("%s:%s\n", name, k)

		//
		// Update Dagger
		//
		dagId, _ := ctx.State().Get("dagger")
		dag, _ := vegdagger.Get(ctx)
		dir := dag.LoadDirectoryFromID(dagger.DirectoryID(dagId.(string)))
		dir = dir.WithNewFile(input.Path, input.Content, dagger.DirectoryWithNewFileOpts{})
		newId, err := dir.ID(ctx)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		//
		// update state
		//
		err = ctx.State().Set(k, input.Content)
		if err != nil {
			return filesysError(input.Path, err), nil
		}
		err = ctx.State().Set("dagger", string(newId))
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

func FilesysDel(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input FilesysPathArgs) (FilesysResult, error) {
		// calculate our real key
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Path)
		fmt.Printf("%s:%s\n", name, k)

		//
		// Update Dagger
		//
		dagId, _ := ctx.State().Get("dagger")
		dag, _ := vegdagger.Get(ctx)
		dir := dag.LoadDirectoryFromID(dagger.DirectoryID(dagId.(string)))
		dir = dir.WithoutFile(input.Path)
		newId, err := dir.ID(ctx)
		if err != nil {
			return filesysError(input.Path, err), nil
		}

		//
		// Add to State ("delete", update)
		//
		err = ctx.State().Set(k, nil)
		if err != nil {
			return filesysError(input.Path, err), nil
		}
		err = ctx.State().Set("dagger", string(newId))
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
