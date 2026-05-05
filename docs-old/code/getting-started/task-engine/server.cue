package examples

import "strings"

server: {
  @flow(server)

  // task: get auth from external command
  gcp: {
    @task(os.Exec)
    cmd: ["gcloud", "auth", "print-access-token"]
    stdout: string
    key:    strings.TrimSpace(stdout)
  }

  run: {
    @task(api.Serve)

    port: "8080"
    routes: {

      // simple hello route
      "/hello": {
        method: "GET"
        resp: {
          status: 200
          body: "hallo chat!"
        }
      }

      // echo request object back as json
      "/jsonecho": {
        method: ["GET", "POST"]
        req: body: {}
        resp: json: req
      }

      // our gemini call warpped in an API
      "/chat": {
        @flow()

        method: "POST"

        // input schema
        req: {
          body: {
            msg: string
          }
        }

        // task: api call via reusable task
        call: _gemini & {
          apikey: gcp.key
          msg: req.body.msg
        }

        // the response to user
        resp: {
          status: 200
          body: call.final.text
        }
      }

    }
  }
}