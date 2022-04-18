package os

import (
  "sort"

	"cuelang.org/go/cue"
  "github.com/mattn/go-zglob"

  hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Glob struct {}

func NewGlob(val cue.Value) (hofcontext.Runner, error) {
  return &FileLock{}, nil
}

func (T *Glob) Run(ctx *hofcontext.Context) (interface{}, error) {

	val := ctx.Value

  patterns, err := extractGlobConfig(ctx, val)
  if err != nil {
    return nil, err
  }

  filepaths, err := filesFromGlobs(patterns)
  if err != nil {
    return nil, err
  }

  return map[string]interface{}{ "filepaths": filepaths}, nil
}

func extractGlobConfig(ctx *hofcontext.Context, val cue.Value) (patterns []string, err error) {
  ctx.CUELock.Lock()
  defer ctx.CUELock.Unlock()

  ps := val.LookupPath(cue.ParsePath("globs"))
  if ps.Err() != nil {
    return nil, ps.Err() 
  }

  iter, err := ps.List()
  if err != nil {
    return nil, err
  }

  for iter.Next() {
    gv := iter.Value()
    if gv.Err() != nil {
      return nil, gv.Err()
    }
    gs, err := gv.String()
    if err != nil {
      return nil, err
    }

    patterns = append(patterns, gs) 
  }

  return patterns, nil 
}

func filesFromGlobs(patterns []string) ([]string, error) {
  // get glob matches
  files := []string{}
  for _, pattern := range patterns {
    matches, err := zglob.Glob(pattern)
    if err != nil {
      return nil, err
    }
    files = append(files, matches...)
  }

  // make unique
  keys := make(map[string]bool)
  unique := make([]string, 0, len(files))	
  for _, file := range files {
      if _, value := keys[file]; !value {
          keys[file] = true
          unique = append(unique, file)
      }
  }    
 
  // also sort
  sort.Strings(unique)
  return unique, nil
}
