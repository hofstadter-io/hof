package db

import "encoding/json"

tasks: {
	@flow(postgres)

	_conn: postgres: "postgres://hof:hof@localhost/hof?sslmode=disable"

	create: {
		@task(db.Call)

		conn: _conn

		exec: #"""
			CREATE TABLE IF NOT EXISTS people (
				name TEXT,
				country TEXT,
				region TEXT,
				occupation TEXT,
				age TEXT
			)
			"""#

		results: _
	}

	insert: {
		@task(db.Call)
		dep:  create
		conn: _conn

		_exec: "INSERT INTO people VALUES ($1, $2, $3, $4, $5)"
		stmts: [ for i in [1, 2, 3, 4, 5] {
			exec: _exec
			args: [ for c in ["n", "c", "r", "o", "a"] {"\(c)\(i)"}]
		}]

		results: _
	}

	query: {
		@task(db.Call)
		dep:  insert
		conn: _conn

		query: "SELECT * FROM people"

		results: _
	}

	destroy: {
		@task(db.Call)
		dep:  query.results
		conn: _conn
		exec: "DROP TABLE people"
	}

	result: {
		@task(os.Stdout)
		dep: [query, destroy]
		text: json.Indent(json.Marshal(query.results), "", "  ") + "\n"
	}

	done: {
		@task(os.Stdout)
		$dep: result
		text: "done\n"
	}
}
