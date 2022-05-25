package db

import (
	"database/sql"
	"fmt"

	"cuelang.org/go/cue"

	_ "github.com/mattn/go-sqlite3"
)

func handleSQLiteExec(dbname, query string, args []interface{}) (string, error) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return "", err
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return "", err
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		return "", err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return "", err
	}

	return fmt.Sprint(affect), nil
}

func handleSQLiteQuery(dbname, query string, args []interface{}) (*sql.Rows, error) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, err
	}

	return db.Query(query)
}

func handleSQLiteStmts(dbname string, stmts cue.Value, args []interface{}) (cue.Value, error) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return stmts, err
	}

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
