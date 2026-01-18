package env

#Method: *"GET" | "POST" | "PUT" | "DELETE" | "OPTIONS" | "HEAD" | "CONNECT" | "TRACE" | "PATCH"

#Route: Ref & {
  // normal http / rest stuff
  // webhook / other validation

  path: string
  method: #Method


  input: {
    url:    string

    headers: [string]: string
    query: [string]:   string

    body: bytes | string | *{} // assumed json body if object
  }

  // some #Thing that get's Sync/Export/Etc...
  vegOp: "SYNC" | "EXPORT" | "CMD"
  handler: _


  routes: [...#Route]
}

