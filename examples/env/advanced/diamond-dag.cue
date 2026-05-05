@experiment(aliasv2)
package advanced

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

diamond: {
  _tmpl: {
    let tmpl = self
    name: string
    port: int | *8080
    _deps: [..._]

		data: env.#Cache & {
      @env()
      #hof: id: "diamond-\(tmpl.name)-data"
      #hof: metadata: name: #hof.id
      name: #hof.id
    }
    ctr: env.#Container & {
      @env()
      #hof: id: "diamond-\(tmpl.name)-ctr"
      #hof: metadata: name: #hof.id

      from: "cgr.dev/chainguard/go"
      steps: [
        env.EnvVars & { FOO: "BAR"},
        env.Mount & { path: "/data", source: tmpl.data },
        env.File & { path: "main.go", content: _diamondCode },
        env.Entrypoint & { args: ["go", "run", "main.go", tmpl.name]},
        for _, dep in tmpl._deps { env.BindService & { service: dep }},
      ]
    }
    svc: env.#Service & {
      @env()
      #hof: id: "diamond-\(tmpl.name)-svc"
      #hof: metadata: name: #hof.id

      hostname: "diamond-\(tmpl.name)"
      source: ctr
      ports: [{port: 8080, frontend: tmpl.port}] 
    }
  }

  // example
  a: _tmpl & { name: "a", port: 8080 }
  b: _tmpl & { name: "b", port: 8080, _deps: [a.svc]}
  c: _tmpl & { name: "c", port: 8080, _deps: [a.svc] }
  d: _tmpl & { name: "d", port: 8080, _deps: [b.svc,c.svc] }
  e: _tmpl & { name: "e", port: 8080, _deps: [a.svc,b.svc,d.svc] }

  // plc
  p_db: _tmpl & { name: "p_db", port: 5432 }
  p: _tmpl & { name: "p", port: 8080, _deps: [p_db.svc] }

  // relay
  r_db: _tmpl & { name: "r_db", port: 5432 }
  r: _tmpl & { name: "r", port: 8080, _deps: [r_db.svc, p.svc] }

  // jetstream
  j: _tmpl & { name: "j", port: 8081, _deps: [p.svc,r.svc] }
  // PDS
  P: _tmpl & { name: "P", port: 8080, _deps: [p.svc,r.svc] }
}

_diamondCode: """
package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello: %v %s\\n", os.Args, r.URL.Path)
	})

	http.ListenAndServe(":8080", nil)
}
"""