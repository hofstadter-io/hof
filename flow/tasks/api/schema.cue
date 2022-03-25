package api

Method: *"GET" | "POST" | "PUT" | "DELETE" | "OPTIONS" | "HEAD" | "CONNECT" | "TRACE" | "PATCH"

Call: {
  @task(api.Call)
  $task: "api.Call"

  req: {
    method: Method
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
 
  port: string
  quitMailbox: string

  routes: [...{
    // @flow() is needed to run sub-tasks per request, which is more typical
    // you can omit if you only need to reshape the data with CUE code

    // filled by hof/flow on each request
    req: {
      method: Method
      url:    string

      headers: [string]: string
      query: [string]:   string

      body: bytes | string | *{} // assumed json body if object

    }

    // any tasks you may need to convert the req -> resp
    // these will be run after the `req` fields has been filled

    // you construct the resp value which is sent back to the client
    // (todo, make this include headers, code, etc
    // for now, this is a value which will be turned into a JSON body for the response
    resp: {
      status: int

      // one of, if none -> NoContent
      json?: {}
      html?: string
      body?: string
    }
    
  }]
}
