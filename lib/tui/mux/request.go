package mux

type Request struct {
	Path    string
	Queries map[string][]string

	Context map[string]interface{}
}
