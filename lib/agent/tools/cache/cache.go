package cache

import (
	"fmt"
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

func CacheRemove(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheRemoveArgs) (CacheResult, error) {
		// calculate our real key
		f := fmt.Sprintf("files:%s:%s", ctx.AgentName(), input.Key)
		fmt.Printf("%s:%s\n", name, f)
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Key)
		fmt.Printf("%s:%s\n", name, k)

		//
		// Add nil to State ("delete", update)
		//
		err := ctx.State().Set(f, nil)
		if err != nil {
			return cacheError(input.Key, err), err
		}
		err = ctx.State().Set(k, nil)
		if err != nil {
			return cacheError(input.Key, err), err
		}

		// return status result
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

func CacheWrite(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheWriteArgs) (CacheResult, error) {
		// calculate our real key
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Key)
		fmt.Printf("%s:%s\n", name, k)

		//
		// update state
		//
		err := ctx.State().Set(k, input.Value)
		if err != nil {
			return cacheError(input.Key, err), err
		}

		// return status result
		return CacheResult{Status: "ok", Key: input.Key}, nil
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

func CacheEdit(name, description string) (tool.Tool, error) {
	handler := func(ctx tool.Context, input CacheEditArgs) (CacheResult, error) {
		k := fmt.Sprintf("cache:%s:%s", ctx.AgentName(), input.Key)
		fmt.Printf("%s:%s\n", name, k)

		//
		// Get from State
		//
		val, err := ctx.State().Get(input.Key)
		if err != nil {
			return cacheError(input.Key, fmt.Errorf("unknown key in statue: %w", err)), nil
		}
		content := val.(string)

		//
		// Check and Replace content
		//
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

		// Add to State
		err = ctx.State().Set(k, next)
		if err != nil {
			return cacheError(input.Key, err), nil
		}

		// return status result
		return CacheResult{Key: input.Key, Status: "ok"}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        name,
		Description: strings.TrimSpace(description),
	}, handler)
}
