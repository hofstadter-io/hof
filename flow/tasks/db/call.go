package db

import (
	"fmt"

	"cuelang.org/go/cue"

  hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Call struct {
  // BaseTask
}

func NewCall(val cue.Value) (hofcontext.Runner, error) {
  return &Call{}, nil
}

func (T *Call) Run(ctx *hofcontext.Context) (interface{}, error) {
  // todo, check failure modes, fill, not return error?
  // (in all tasks)
  // does a failed stmt fail the proper way, especially when transaction?

	v := ctx.Value

  out, err := handleQuery(ctx, v)
  if err != nil {
    return nil, err
  }

  ctx.CUELock.Lock()
  defer ctx.CUELock.Unlock()
	return v.FillPath(cue.ParsePath("results"), out), nil
}

func handleQuery(ctx *hofcontext.Context, val cue.Value) (interface{}, error) {

  var (
    query    cue.Value
    callType string
    dbname   string
    qs       string
    iargs    []interface{}
    err      error
  )


  ferr := func () error {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()

    query = val.LookupPath(cue.ParsePath("query"))
    if query.Exists() && query.Err() == nil {
      callType = "query"
    }
    if callType == "" {
      query = val.LookupPath(cue.ParsePath("exec"))
      if query.Exists() && query.Err() == nil {
        callType = "exec"
      }
    }
    if callType == "" {
      query = val.LookupPath(cue.ParsePath("stmts"))
      if query.Exists() && query.Err() == nil {
        callType = "stmts"
      }
    }

    var args []string
    av := val.LookupPath(cue.ParsePath("args"))
    if av.Exists() {
      err = av.Decode(&args)
      if err != nil {
        return fmt.Errorf("while decoding 'args' at %v", err)
      }
    }

    // args for the database call
    for _, a := range args {
      iargs = append(iargs, a)
    }

    conn := val.LookupPath(cue.ParsePath("conn"))
    if !conn.Exists() {
      return fmt.Errorf("field 'conn' is required on db.Query at %q", val.Path())
    }

    // what type of database is this?
    // we shouldn't look here, migrate to a lookup table

    iter, err := conn.Fields()
    if err != nil {
      return fmt.Errorf("in field 'conn' at %v", err)
    }


    for iter.Next() {
      sel := iter.Selector().String()
      switch sel {
      case "sqlite":
        dbname, err = iter.Value().String()
        if err != nil {
          return err
        }
        if callType != "stmts" {
          qs, err = query.String()
          if err != nil {
            return fmt.Errorf("in field 'query' at %v", err)
          }
        }

      }
    }
    return nil
  }()
  if ferr != nil {
    return nil, ferr
  }

  switch callType {
  case "query":

    rows, err := handleSQLiteQuery(dbname, qs, iargs)
    if err != nil {
      return nil, fmt.Errorf("error during query %v", err)
    }

    jstr, err := scanRowToJson(rows)
    if err != nil {
      return nil, fmt.Errorf("error during scan %v", err)
    }
    return val.Context().CompileBytes(jstr), nil

  case "exec":
    out, err := handleSQLiteExec(dbname, qs, iargs)
    if err != nil {
      return nil, fmt.Errorf("error during exec %v", err)
    }
    return out, nil

  case "stmts":
    out, err := handleSQLiteStmts(dbname, query, iargs)
    if err != nil {
      return nil, fmt.Errorf("error during query %v", err)
    }
    return out, nil
  }

	return "", fmt.Errorf("no supported conn types found in db.Query %q", val.Path())
}
