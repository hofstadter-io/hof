#HttpMethod: "OPTIONS" | "HEAD" | "GET" | "POST" | "PATCH" | "PUT" | "DELETE" | "CONNECT" | "TRACE"

#Server: {
	// ...

	Routes: #Routes
}

#Routes: [...#Route] | *[]
#Route: {
	Name:   string
	Path:   string
	Method: #HttpMethod

	// Route and Query params
	Params: [...string] | *[]
	Query:  [...string] | *[]

	// Fields which allow the user to write
	// handler bodies directly in CUE
	Body?:   string
	Imports: [...string] | *[]

	// Allows subroutes for routes
	Routes: [...#Route]
}
