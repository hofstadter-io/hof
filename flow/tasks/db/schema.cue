package db

// need to figure out how to separate out the connection
// so it can be create once and reused by queries
// possibly how to pool connections as well

// only SQLite supported currently

// Call a database
Call: {
  @task(db.Call)
  $task: "db.Call"

  // db connection
  conn: {
    sqlite: string // db name
  }

  // args to Call
  args:  [..._]

  // Use only one of [query,exec,stmts] 
  query: string
  exec:  string

  stmts: [...{
    // Use only one of [query,exec,stmts] 
    query: string
    exec:  string
    // args to statement, merged with top-level
    args:  string 
  }]

  results: _
}

