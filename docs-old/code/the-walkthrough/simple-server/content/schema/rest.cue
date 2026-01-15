package schema

HttpMethod: "OPTIONS" | "HEAD" | "GET" | "POST" | "PATCH" | "PUT" | "DELETE" | "CONNECT" | "TRACE"

Server: {
	// schemas typically have a name field
	//   the user sets this for each instance
	Name: string

	// Some more common "optional" fields
	// we use defaults rather than CUE optional syntax
	Description: string | *""
	Help:        string | *""

	// the REST routes
	Routes: Routes
}
