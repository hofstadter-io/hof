@flow()

nested: {
  @task(nest)

  call: {
    @task(api.Call)
    req: {
      host: "https://postman-echo.com"
      method: "GET"
      path: "/get"
      query: {
        cow: "moo"
      }
    }
    resp: string
  }
  tmp: { 
    @task(noop)
    resp: call.resp
  }
}

out: { text: nested.tmp.resp + "\n" } @task(os.Stdout)
