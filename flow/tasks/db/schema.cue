package db

Query: {
  @task(db.Query)
  conn: {...}
  query: string
  args: [..._]
  results: _
}

