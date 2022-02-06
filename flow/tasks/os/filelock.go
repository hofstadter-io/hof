package os

import (
  "fmt"
  "time"

	"cuelang.org/go/cue"
  "github.com/gofrs/flock"

  "github.com/hofstadter-io/hof/flow/context"
)

func init() {
  context.Register("os.FileLock", NewFileLock)
  context.Register("os.FileUnlock", NewFileUnlock)
}

type FileLock struct {}

func NewFileLock(val cue.Value) (context.Runner, error) {
  return &FileLock{}, nil
}

func (T *FileLock) Run(ctx *context.Context) (interface{}, error) {

	val := ctx.Value

  fn, rw, retry, err := extractConfig(ctx, val)
  if err != nil {
    return nil, err
  }

  kn := "hof-filelock."+fn

  // check store for existing filelock

  var lock *flock.Flock
  l, ok := ctx.ValStore.Load(kn)
  if !ok {
    lock = flock.New(fn)
  } else {
    lock = l.(*flock.Flock)
  }

  ctx.ValStore.Store(kn, lock)

  if retry == 0 {
    if rw {
      err = lock.Lock()
    } else {
      err = lock.RLock()
    }
  } else {
    if rw {
      _, err = lock.TryLockContext(ctx.Context, retry)
    } else {
      _, err = lock.TryRLockContext(ctx.Context, retry)
    }
  }

	return nil, err
}


type FileUnlock struct {}

func NewFileUnlock(val cue.Value) (context.Runner, error) {
  return &FileUnlock{}, nil
}

func (T *FileUnlock) Run(ctx *context.Context) (interface{}, error) {

	val := ctx.Value

  fn, _, _, err := extractConfig(ctx, val)
  if err != nil {
    return nil, err
  }

  kn := "hof-filelock."+fn

  l, ok := ctx.ValStore.Load(kn)
  if !ok {
    return nil, fmt.Errorf("unknown filelock: %q", fn)
  }
  lock := l.(*flock.Flock)

  err = lock.Unlock()

	return nil, err
}

func extractConfig(ctx *context.Context, val cue.Value) (fn string, rw bool, retry time.Duration, err error) {
  ctx.CUELock.Lock()
  defer ctx.CUELock.Unlock()

  f := val.LookupPath(cue.ParsePath("filename")) 
  if f.Err() != nil {
    return fn, rw, retry, f.Err()
  } else if f.Exists() {
    fn, err = f.String()
    if err != nil {
      return fn, rw, retry, err
    }
  }

  r := val.LookupPath(cue.ParsePath("rw"))
  if r.Exists() {
    rw, err = r.Bool()
    if err != nil {
      return fn, rw, retry, err
    }
  }

  d := val.LookupPath(cue.ParsePath("retry")) 
  if d.Exists() {
    ds, err := d.String()
    if err != nil {
      return fn, rw, retry, err
    }
    retry, err = time.ParseDuration(ds)
    if err != nil {
      return fn, rw, retry, err
    }
  }

  return fn, rw, retry, err
}
