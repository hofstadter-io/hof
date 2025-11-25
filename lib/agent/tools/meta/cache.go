package meta

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

const CacheRemoveTool = "cache_remove"
const CacheRemoveDesc = `
adds or overwrites an entry in your working key/value cache and context
`

type CacheRemoveArgs struct {
	Key string `json:"key"` // path to a directory
}
type CacheRemoveResult struct {
	Key    string `json:"key"`    // path to a directory
	Status string `json:"status"` // "ok" or "error"
	Error  string `json:"error,omitempty"`
}

func NewCacheRemove() (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheRemoveArgs) (CacheRemoveResult, error) {
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Key)
		fmt.Println("CACHE.REMOVE:", k)
		ctx.Actions().StateDelta[k] = nil
		// err := ctx.State().Set(k, nil)
		// if err != nil {
		// 	return CacheRemoveResult{Status: "error", Key: input.Key, Error: err.Error()}, err
		// }
		return CacheRemoveResult{Status: "ok", Key: input.Key}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        CacheRemoveTool,
		Description: strings.TrimSpace(CacheRemoveDesc),
	}, handler)
}

const CacheWriteTool = "cache_write"
const CacheWriteDesc = `
adds or overwrites an entry in your working key/value cache and context
`

type CacheWriteArgs struct {
	Key   string `json:"key"`   // path to a directory
	Value string `json:"value"` // path to a directory
}
type CacheWriteResult struct {
	Key    string `json:"key"`    // path to a directory
	Status string `json:"status"` // "ok" or "error"
	Error  string `json:"error,omitempty"`
}

func NewCacheWrite() (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheWriteArgs) (CacheWriteResult, error) {
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Key)
		fmt.Println("CACHE.WRITE:", k)

		// err := ctx.State().Set(k, input.Value)
		// if err != nil {
		// 	return CacheWriteResult{Status: "error", Key: input.Key, Error: err.Error()}, err
		// }

		ctx.Actions().StateDelta[k] = input.Value

		return CacheWriteResult{Status: "ok", Key: input.Key}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        CacheWriteTool,
		Description: strings.TrimSpace(CacheWriteDesc),
	}, handler)
}

const CacheFileTool = "cache_file"
const CacheFileDesc = `
read the file at $path and puts the contents in your working key/value cache and context
`

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

func NewCacheFile() (tool.Tool, error) {
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
		ctx.State().Set(k, string(c))
		// return the result
		return CacheFileResult{Path: input.Path, Status: "ok"}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        CacheFileTool,
		Description: strings.TrimSpace(CacheFileDesc),
	}, handler)
}

const CacheDirTool = "cache_dir"
const CacheDirDesc = `
lists a directory and stores it in your working key/value cache and context
`

type CacheDirArgs struct {
	Path string `json:"path"` // path to a directory
}
type CacheDirResult struct {
	Path   string `json:"path"`   // path to a directory
	Status string `json:"status"` // "ok" or "error"
	Error  string `json:"error,omitempty"`
}

func NewCacheDir() (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheDirArgs) (CacheDirResult, error) {
		state := maps.Collect(ctx.State().All())
		fmt.Println("read_dir:", input.Path, state)
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

		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Path)
		ctx.State().Set(k, b.String())
		return CacheDirResult{Status: "ok", Path: input.Path}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        CacheDirTool,
		Description: strings.TrimSpace(CacheDirDesc),
	}, handler)
}
