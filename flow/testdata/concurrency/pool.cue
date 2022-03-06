package concurrency

import "list"
import "strconv"

// demonstrates limiting a task to
// at most N concurrent processors
poolExample: {
  @flow(pool/exec)

  maxS: string | *"5" @tag(max)
  max: strconv.Atoi(maxS) 

  init: {
    @task(noop)
    @pool(api,2) // this is a problem, can't pass a dynamic value
  }

  for i,_ in list.Range(0,max,1) {
    "task-\(i)": {
      // @pool(api)
      @task(nest)
      call: {
        @task(api.Call)
        req: {
          host: "https://postman-echo.com"
          method: "GET"
          path: "/get"
          query: {
            cow: "moo"
            task: "\(i)"
          }
        }
        resp: string
      }
      out: { text: call.resp + "\n"} @task(os.Stdout)
      wait: { duration: "1s", dep: [call,out] } @task(os.Sleep)
    }

  }
}

