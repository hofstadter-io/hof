package schema

HttpMethod: "OPTIONS" | "HEAD" | "GET" | "POST" | "PATCH" | "PUT" | "DELETE" | "CONNECT" | "TRACE"

Server: {
	// Most schemas have a name field
	Name: string

	// Some more common "optional" fields
	// we use defaults rather than CUE optional syntax
	Description: string | *""
	Help:        string | *""

	// the REST routes
	Routes: Routes
}
