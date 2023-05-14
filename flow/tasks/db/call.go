package db

import (
	"database/sql"
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

func (T *Call) Run(ctx *hofcontext.Context) (any, error) {
	// todo, check failure modes, fill, not return error?
	// (in all tasks)
	// does a failed stmt fail the proper way, especially when transaction?

	v := ctx.Value

	out, err := handleQueryCall(ctx, v)
	if err != nil {
		return nil, err
	}

	ctx.CUELock.Lock()
	defer ctx.CUELock.Unlock()
	return v.FillPath(cue.ParsePath("results"), out), nil
}

func handleQueryCall(ctx *hofcontext.Context, val cue.Value) (any, error) {

	var (
		query    cue.Value
		callType string
		dbtype   string
		dbname   string
		qs       string
		iargs    []interface{}
		err      error
	)

	// func to procect access to CUE value with lock
	// while CUE is not concurrency safe
	ferr := func() error {
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
				dbtype = "sqlite"
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

			case "postgres":
				dbtype = "postgres"
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
			default:
				return fmt.Errorf("unknown db conn type %q", sel)
			}

		}
		return nil
	}()
	if ferr != nil {
		return nil, ferr
	}

	switch callType {
	case "query":

		rows, err := handleQueryDB(dbtype, dbname, qs, iargs)
		if err != nil {
			return nil, fmt.Errorf("error during query %v", err)
		}

		jstr, err := scanRowToJson(rows)
		if err != nil {
			return nil, fmt.Errorf("error during scan %v", err)
		}
		return val.Context().CompileBytes(jstr), nil

	case "exec":
		out, err := handleExecDB(dbtype, dbname, qs, iargs)
		if err != nil {
			return nil, fmt.Errorf("error during exec %v", err)
		}
		return out, nil

	case "stmts":
		out, err := handleStmtsDB(dbtype, dbname, query, iargs)
		if err != nil {
			return nil, fmt.Errorf("error during query %v", err)
		}
		return out, nil
	}

	return "", fmt.Errorf("no supported conn types found in db.Query %q", val.Path())
}

func handleExecDB(dbtype, dbname, query string, args []interface{}) (string, error) {
	switch dbtype {
	case "sqlite":
		return handleSQLiteExec(dbname, query, args)
	case "postgres":
		return handlePostgresExec(dbname, query, args)
	default:
		return "", fmt.Errorf("unknown db type: %q", dbtype)
	}
}

func handleExec(db *sql.DB, query string, args []interface{}) (string, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println("got here 1")
		return "", err
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		fmt.Println("got here 2")
		return "", err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return "", err
	}

	return fmt.Sprint(affect), nil
}

func handleQueryDB(dbtype, dbname, query string, args []interface{}) (*sql.Rows, error) {

	switch dbtype {
	case "sqlite":
		return handleSQLiteQuery(dbname, query, args)
	case "postgres":
		return handlePostgresQuery(dbname, query, args)
	default:
		return nil, fmt.Errorf("unknown db type: %q", dbtype)
	}
}

func handleQuery(db *sql.DB, query string, args []interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func handleStmtsDB(dbtype, dbname string, stmts cue.Value, args []interface{}) (cue.Value, error) {
	switch dbtype {
	case "sqlite":
		return handleSQLiteStmts(dbname, stmts, args)
	case "postgres":
		return handlePostgresStmts(dbname, stmts, args)
	default:
		return stmts, fmt.Errorf("unknown db type: %q", dbtype)
	}
}

func handleStmts(db *sql.DB, stmts cue.Value, args []interface{}) (cue.Value, error) {
	iter, err := stmts.List()
	if err != nil {
		return stmts, err
	}

	results := []cue.Value{}
	for iter.Next() {
		val := iter.Value()
		sel := iter.Selector()
		callType := ""

		query := val.LookupPath(cue.ParsePath("query"))
		if query.Exists() && query.Err() == nil {
			rows, err := db.Query(query.String())
			if err != nil {
				return stmts, fmt.Errorf("error during scan %v", err)
			}
			jstr, err := scanRowToJson(rows)
			if err != nil {
				return stmts, fmt.Errorf("error during scan %v", err)
			}
			r := val.Context().CompileBytes(jstr)
			val = val.FillPath(cue.ParsePath("results"), r)
			results = append(results, val)
			continue
		}

		query = val.LookupPath(cue.ParsePath("exec"))
		if query.Exists() && query.Err() == nil {
			qs, err := query.String()
			stmt, err := db.Prepare(qs)
			if err != nil {
				return stmts, err
			}

			// handle local args
			var la []string
			av := val.LookupPath(cue.ParsePath("args"))
			if av.Exists() {
				err = av.Decode(&la)
				if err != nil {
					return stmts, fmt.Errorf("while decoding 'args' at %v", err)
				}
			}

			ia := []interface{}{}
			for _, a := range la {
				ia = append(ia, a)
			}

			res, err := stmt.Exec(ia...)
			if err != nil {
				return stmts, err
			}

			affect, err := res.RowsAffected()
			if err != nil {
				return stmts, err
			}

			val = val.FillPath(cue.ParsePath("results"), fmt.Sprint(affect))
			results = append(results, val)
			continue
		}

		fmt.Println(sel, callType)

		// do db calls
		// fill val

		results = append(results, val)
	}


	return stmts.Context().NewList(results...), nil
}

