package api

Call: {
  @task(api.Call)
  $task: "api.Call"

  req: {
    method: *"GET" | "POST" | "PUT" | "DELETE" | "OPTIONS" | "HEAD" | "CONNECT" | "TRACE" | "PATCH"
    host:   string
    path:   string | *""
    auth?:  string
    headers?: [string]: string
    query?: [string]:   string
    form?: [string]:   string
    data?:    string | {...}
    timeout?: string
    // curl?: string
    retry?: {
      count: int | *3
      timer: string | *"6s"
      codes: [...int]
    }
  }

  // filled by task
  resp: {
		status:     string
		statusCode: int

		body: *{} | bytes | string
		header: [string]:  string | [...string]
		trailer: [string]:  string | [...string]
  }

}

Serve: {
  @task(api.Serve)
  $task: "api.Serve"
}
