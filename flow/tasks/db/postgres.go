package db

import (
	"database/sql"

	"cuelang.org/go/cue"

	_ "github.com/lib/pq"
)

func handlePostgresExec(dbname, query string, args []interface{}) (string, error) {
	db, err := sql.Open("postgres", dbname)
	if err != nil {
		return "", err
	}
	return handleExec(db, query, args)
}

func handlePostgresQuery(dbname, query string, args []interface{}) (*sql.Rows, error) {
	db, err := sql.Open("postgres", dbname)
	if err != nil {
		return nil, err
	}
	return handleQuery(db, query, args)
}
func handlePostgresStmts(dbname string, stmts cue.Value, args []interface{}) (cue.Value, error) {
	db, err := sql.Open("postgres", dbname)
	if err != nil {
		return stmts, err
	}
	return handleStmts(db, stmts, args)
}
