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
type CacheResult struct {
	Key    string `json:"key"`    // path to a directory
	Status string `json:"status"` // "ok" or "error"
	Error  string `json:"error,omitempty"`
}

func cacheError(key string, err error) CacheResult {
	fmt.Println("ERROR:", key, err)
	return CacheResult{Key: key, Status: "error", Error: err.Error()}
}

func NewCacheRemove(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheRemoveArgs) (CacheResult, error) {
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Key)
		// TODO, do we also need to set StateDelta so the next function sees it?
		// ctx.Actions().StateDelta[k] = nil
		fmt.Println("cache_remove:", k)
		err := ctx.State().Set(k, nil)
		if err != nil {
			return cacheError(input.Key, err), err
		}
		return CacheResult{Status: "ok", Key: input.Key}, nil
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

func NewCacheWrite(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheWriteArgs) (CacheResult, error) {
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Key)

		// ctx.Actions().StateDelta[k] = input.Value
		fmt.Println("cache_write:", k, len(input.Value))
		err := ctx.State().Set(k, input.Value)
		if err != nil {
			return cacheError(input.Key, err), err
		}

		return CacheResult{Status: "ok", Key: input.Key}, nil
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

func NewCacheFile(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheFileArgs) (CacheResult, error) {
		// need to cascade where we look, from closest to outer most
		// 1. state
		// 2. artifacts
		// 3. filesystem (host vs ephemeral) (tied to terminal access/env)
		// same for most tools

		state := maps.Collect(ctx.State().All())
		// fmt.Println("cache_file:", input.Path, state)

		// get the env from state
		env, ok := state["env"].(map[string]any)
		if !ok {
			err := fmt.Errorf("while finding env state")
			return cacheError(input.Path, err), err
		}
		// get the workspace dir from env
		wsDir, ok := env["workspaceDir"].(string)
		if !ok {
			err := fmt.Errorf("while finding workspace directory")
			return cacheError(input.Path, err), err
		}

		// read file relative to workspace dir
		c, err := os.ReadFile(filepath.Join(wsDir, input.Path))
		if err != nil {
			return cacheError(input.Path, err), err
		}
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Path)
		// ctx.Actions().StateDelta[k] = string(c)

		fmt.Println("cache_file:", k)
		err = ctx.State().Set(k, string(c))
		if err != nil {
			return cacheError(input.Path, err), err
		}
		// return the result
		return CacheResult{Key: input.Path, Status: "ok"}, nil
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
	handler := func(ctx tool.Context, input CacheDirArgs) (CacheResult, error) {
		state := maps.Collect(ctx.State().All())
		fmt.Println("cache_dir", input.Path, state)
		env, ok := state["env"].(map[string]any)
		if !ok {
			err := fmt.Errorf("failed to get env")
			return cacheError(input.Path, err), err
		}
		wsDir, ok := env["workspaceDir"].(string)
		if !ok {
			err := fmt.Errorf("failed to get wsDir")
			return cacheError(input.Path, err), err
		}

		entries, err := os.ReadDir(filepath.Join(wsDir, input.Path))
		if err != nil {
			return cacheError(input.Path, err), err
		}

		b := new(strings.Builder)
		for _, e := range entries {
			fmt.Fprintln(b, e.Name())
		}

		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Path)

		// ctx.Actions().StateDelta[k] = b.String()
		fmt.Println("cache_dir:", k)
		err = ctx.State().Set(k, b.String())
		if err != nil {
			return cacheError(input.Path, err), err
		}
		return CacheResult{Status: "ok", Key: input.Path}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}

type CacheEditArgs struct {
	Key string `json:"key"` // path to a file
	Old string `json:"old_string"`
	New string `json:"new_string"`
	Exp int    `json:"expected replacements"` // defaults to 1 if not set
}

func NewCacheEdit(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheEditArgs) (CacheResult, error) {
		// need to cascade where we look, from closest to outer most
		// 1. state
		// 2. filesystem (host vs ephemeral) (tied to terminal access/env)
		// same for most tools
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Key)

		curr, err := ctx.State().Get(k)
		if err != nil {
			return cacheError(input.Key, err), nil
		}

		content, ok := curr.(string)
		if !ok {
			err = fmt.Errorf("key %q was not found in the cache", input.Key)
			return cacheError(input.Key, err), nil
		}

		count := input.Exp
		if count == 0 {
			count = 1
		}
		found := strings.Count(content, input.Old)
		if found != count {
			err = fmt.Errorf("while editing %q, expected %d matches, but found %d", input.Key, count, found)
			return cacheError(input.Key, err), nil
		}

		next := strings.Replace(content, input.Old, input.New, count)

		// ctx.Actions().StateDelta[k] = next
		err = ctx.State().Set(k, next)
		if err != nil {
			return cacheError(input.Key, err), nil
		}
		// return the result
		return CacheResult{Key: input.Key, Status: "ok"}, nil
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
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Path)

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
			err = fmt.Errorf("Error:\n%w\n\nOutput:\n%s\n\n", err, string(out))
			return cacheGrepError(input.Path, err), err
		}

		// maybe count '\n' and split instead? Ideally we could token budget
		if len(out) > 4000 {
			out = out[:4000]
			out = append(out, []byte("\n...\noutput is too long and has been truncated")...)
		}

		// ctx.Actions().StateDelta[k] = string(out)
		err = ctx.State().Set(k, string(out))
		if err != nil {
			return cacheGrepError(input.Path, err), err
		}
		// return the result
		return CacheGrepResult{Path: input.Path, Status: "ok"}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}
