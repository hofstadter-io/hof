package meta

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type CacheRemoveArgs struct {
	Key string `json:"key"` // path to a directory
}
type CacheRemoveResult struct {
	Key    string `json:"key"`    // path to a directory
	Status string `json:"status"` // "ok" or "error"
	Error  string `json:"error,omitempty"`
}

func NewCacheRemove(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheRemoveArgs) (CacheRemoveResult, error) {
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Key)
		ctx.Actions().StateDelta[k] = nil
		// TODO, do we also need to set state so the next function sees it?
		// err := ctx.State().Set(k, nil)
		// if err != nil {
		// 	return CacheRemoveResult{Status: "error", Key: input.Key, Error: err.Error()}, err
		// }
		return CacheRemoveResult{Status: "ok", Key: input.Key}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

type CacheWriteArgs struct {
	Key   string `json:"key"`   // path to a directory
	Value string `json:"value"` // path to a directory
}
type CacheWriteResult struct {
	Key    string `json:"key"`    // path to a directory
	Status string `json:"status"` // "ok" or "error"
	Error  string `json:"error,omitempty"`
}

func NewCacheWrite(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheWriteArgs) (CacheWriteResult, error) {
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Key)

		ctx.Actions().StateDelta[k] = input.Value

		return CacheWriteResult{Status: "ok", Key: input.Key}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

type CacheFileArgs struct {
	Path string `json:"path"` // path to a file
}
type CacheFileResult struct {
	Path   string `json:"path"`            // path to a file
	Status string `json:"status"`          // "ok" or "error"
	Error  string `json:"error,omitempty"` // an error message
}

func cacheFileError(path string, err error) CacheFileResult {
	return CacheFileResult{Path: path, Status: "error", Error: err.Error()}
}

func NewCacheFile(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheFileArgs) (CacheFileResult, error) {
		// need to cascade where we look, from closest to outer most
		// 1. state
		// 2. artifacts
		// 3. filesystem (host vs ephemeral) (tied to terminal access/env)
		// same for most tools

		state := maps.Collect(ctx.State().All())
		fmt.Println("cache_file:", input.Path, state)

		// get the env from state
		env, ok := state["env"].(map[string]any)
		if !ok {
			err := fmt.Errorf("while finding env state")
			return cacheFileError(input.Path, err), err
		}
		// get the workspace dir from env
		wsDir, ok := env["workspaceDir"].(string)
		if !ok {
			err := fmt.Errorf("while finding workspace directory")
			return cacheFileError(input.Path, err), err
		}

		// read file relative to workspace dir
		c, err := os.ReadFile(filepath.Join(wsDir, input.Path))
		if err != nil {
			fmt.Println("while reading from filesystem:", err)
			return cacheFileError(input.Path, err), err
		}
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Path)
		ctx.Actions().StateDelta[k] = string(c)
		// ctx.State().Set(k, string(c))
		// return the result
		return CacheFileResult{Path: input.Path, Status: "ok"}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

type CacheDirArgs struct {
	Path string `json:"path"` // path to a directory
}
type CacheDirResult struct {
	Path   string `json:"path"`   // path to a directory
	Status string `json:"status"` // "ok" or "error"
	Error  string `json:"error,omitempty"`
}

func NewCacheDir(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheDirArgs) (CacheDirResult, error) {
		state := maps.Collect(ctx.State().All())
		fmt.Println("cache_dir", input.Path, state)
		env, ok := state["env"].(map[string]any)
		if !ok {
			err := fmt.Errorf("failed to get env")
			return CacheDirResult{Path: input.Path, Status: "error", Error: err.Error()}, err
		}
		wsDir, ok := env["workspaceDir"].(string)
		if !ok {
			err := fmt.Errorf("failed to get wsDir")
			return CacheDirResult{Path: input.Path, Status: "error", Error: err.Error()}, err
		}

		entries, err := os.ReadDir(filepath.Join(wsDir, input.Path))
		if err != nil {
			return CacheDirResult{Path: input.Path, Status: "error", Error: err.Error()}, err
		}

		b := new(strings.Builder)
		for _, e := range entries {
			fmt.Fprintln(b, e.Name())
		}

		k := fmt.Sprintf("cache:%s:dir:%s", ctx.AgentName(), input.Path)

		ctx.Actions().StateDelta[k] = b.String()
		// ctx.State().Set(k, b.String())
		return CacheDirResult{Status: "ok", Path: input.Path}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

type CacheEditArgs struct {
	Path string `json:"path"` // path to a file
	Old  string `json:"old_string"`
	New  string `json:"new_string"`
	Exp  int    `json:"expected replacements"`
}

type CacheEditResult struct {
	Path   string `json:"path"`            // path to a file
	Status string `json:"status"`          // "ok" or "error"
	Error  string `json:"error,omitempty"` // an error message
}

func cacheEditError(path string, err error) CacheEditResult {
	return CacheEditResult{Path: path, Status: "error", Error: err.Error()}
}

func NewCacheEdit(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheEditArgs) (CacheEditResult, error) {
		// need to cascade where we look, from closest to outer most
		// 1. state
		// 2. filesystem (host vs ephemeral) (tied to terminal access/env)
		// same for most tools
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Path)

		curr, err := ctx.State().Get(k)
		if err != nil {
			return cacheEditError(input.Path, err), err
		}

		content, ok := curr.(string)
		if !ok {
			err = fmt.Errorf("%s was not found in the cache", input.Path)
			return cacheEditError(input.Path, err), err
		}

		count := input.Exp
		if count == 0 {
			count = 1
		}
		found := strings.Count(content, input.Old)
		if found != count {
			err = fmt.Errorf("while editing %s, expected %d matches, but found %d", input.Path, count, found)
			return cacheEditError(input.Path, err), err
		}

		next := strings.Replace(content, input.Old, input.New, count)

		ctx.Actions().StateDelta[k] = next
		// return the result
		return CacheEditResult{Path: input.Path, Status: "ok"}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

type CacheGrepResult struct {
	Path   string `json:"path"`            // path to a file
	Status string `json:"status"`          // "ok" or "error"
	Error  string `json:"error,omitempty"` // an error message
}

type CacheGrepArgs struct {
	Path     string `json:"path"`         // base path to grep from
	Regexp   string `json:"regexp"`       // a regular expression to grep for
	Around   int    `json:"lines_around"` // the number of surrounding lines to include, max 8, defaults to 0 (only matching lines)
	MaxDepth int    `json:"max_depth"`    // the recursive level to grep, max 3, defaults to 1
}

func cacheGrepError(path string, err error) CacheGrepResult {
	return CacheGrepResult{Path: path, Status: "error", Error: err.Error()}
}

func NewCacheGrep(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheGrepArgs) (CacheGrepResult, error) {
		// need to cascade where we look, from closest to outer most
		// 1. state
		// 2. filesystem (host vs ephemeral) (tied to terminal access/env)
		// same for most tools
		k := fmt.Sprintf("cache:%s:grep:%s", ctx.AgentName(), input.Path)

		// TODO, is there a stdlib "clamp" function in Go?
		if input.Around < 0 {
			input.Around = 0
		}
		if input.Around > 8 {
			input.Around = 8
		}
		if input.MaxDepth < 1 {
			input.MaxDepth = 1
		}
		if input.MaxDepth > 3 {
			input.MaxDepth = 3
		}

		scriptFmt := `
		rg -Rn -B%d -A%d --sort=path -e '%s' %s
		`
		// -B%s -A%d  // for extra lines before / after
		// limit to a set of globs?

		script := strings.TrimSpace(fmt.Sprintf(scriptFmt, input.Around, input.Around, input.Regexp, input.Path))

		cmd := exec.Command("sh", "-c", script)
		if input.Path != "" {
			cmd.Dir = input.Path
		}
		out, err := cmd.CombinedOutput()
		if err != nil {
			err = fmt.Errorf("Error:\n%s\n\nOutput:\n%s\n\n", err, string(out))
			return cacheGrepError(input.Path, err), err
		}

		// maybe count '\n' and split instead? Ideally we could token budget
		if len(out) > 4000 {
			out = out[:4000]
			out = append(out, []byte("\n...\noutput is too long and has been truncated")...)
		}

		ctx.Actions().StateDelta[k] = string(out)
		// ctx.State().Set(k, string(out))
		// return the result
		return CacheGrepResult{Path: input.Path, Status: "ok"}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}
