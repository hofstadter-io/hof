package examples

import (
	"hof.io/docs/example/schema"
)

// A concrete value of the Server schem
Server: schema.Server & {
	GoModule: "hof.io/docs/example"

	Name:        "Example"
	Description: "An example server"

	Routes: [{
		Name:   "EchoQ"
		Path:   "/echo"
		Method: "GET"
		Query: ["msg"]
		Body: """
			c.String(http.StatusOK, msg)
			"""
	}, {
		Name:   "EchoP"
		Path:   "/echo"
		Method: "GET"
		Params: ["msg"]
		Body: """
			c.String(http.StatusOK, msg)
			"""
	}, {
		Name:   "Hello"
		Path:   "/hello"
		Method: "GET"
		Query: ["msg"]
		Body: """
			if msg == "" {
				msg = "hello world"
			}
			c.String(http.StatusOK, msg)
			"""
	}]

	Prometheus: true
}
