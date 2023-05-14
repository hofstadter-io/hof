package db

import (
	"database/sql"

	"cuelang.org/go/cue"

	_ "modernc.org/sqlite"
)

func handleSQLiteExec(dbname, query string, args []interface{}) (string, error) {
	db, err := sql.Open("sqlite", dbname)
	if err != nil {
		return "", err
	}
	return handleExec(db, query, args)
}

func handleSQLiteQuery(dbname, query string, args []interface{}) (*sql.Rows, error) {
	db, err := sql.Open("sqlite", dbname)
	if err != nil {
		return nil, err
	}

	return handleQuery(db, query, args)
}

func handleSQLiteStmts(dbname string, stmts cue.Value, args []interface{}) (cue.Value, error) {
	db, err := sql.Open("sqlite", dbname)
	if err != nil {
		return stmts, err
	}

	return handleStmts(db, stmts, args)
}
