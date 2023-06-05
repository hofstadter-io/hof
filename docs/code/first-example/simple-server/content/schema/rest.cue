package schema

Server: {
	// Most schemas have a name field
	Name: string

	// Some more common "optional" fields
	// we use defaults rather than CUE optional syntax
	Description: string | *""
	Help:        string | *""
}
