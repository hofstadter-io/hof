package schema

import (
	"strings"

)

#HttpMethod: "OPTIONS" | "HEAD" | "GET" | "POST" | "PATCH" | "PUT" | "DELETE" | "CONNECT" | "TRACE"

#Server: {
	// Most schemas have a name field
	Name: string

	// Some more common "optional" fields
	// we use defaults rather than CUE optional syntax
	Description: string | *""
	Help:        string | *""

	// The server routes
	Routes: #Routes

	// list of file globs to be embedded into the server when built
	// todo, pass this back through the sections if it is not
	StaticFiles: [...#Static]
	// enable prometheus metrics
	Prometheus: bool | *false

	// various casings of the server Name
	serverName:  strings.ToCamel(Name)
	ServerName:  strings.ToTitle(Name)
	SERVER_NAME: strings.ToUpper(Name)
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

#Resources: [string]: #Resource
#Resource: {
	Model:  #Model
	Name:   Model.Name
	Routes: #Routes
}

#DatamodelToResources: {
	Datamodel: #Datamodel
	Resources: #Resources & {
		for n, M in Datamodel.Models {
			"\(n)": {
				Model: M
				Name:  M.Name
				Routes: [{
					Name:   "\(M.Name)Create"
					Path:   ""
					Method: "POST"
				}, {
					Name: "\(M.Name)Read"
					Path: ""
					Params: ["\(strings.ToLower(M.Index))"]
					Method: "GET"
				}, {
					Name:   "\(M.Name)List"
					Path:   ""
					Method: "GET"
				}, {
					Name:   "\(M.Name)Update"
					Path:   ""
					Method: "PATCH"
				}, {
					Name: "\(M.Name)Delete"
					Path: ""
					Params: ["\(strings.ToLower(M.Index))"]
					Method: "DELETE"
				}, ...] // left open so you can add custom routes
			}
		}
	}
}

#Static: {
	Globs: [...string]
	TrimPrefix?: string
	AddPrefix?:  string
}
