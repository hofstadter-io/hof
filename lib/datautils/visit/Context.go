package main

/*
Where's your docs doc?!
*/
type Context struct {
	Obj  interface{} `json:"obj" xml:"obj" yaml:"obj" form:"obj" query:"obj" `
	Path []string    `json:"path" xml:"path" yaml:"path" form:"path" query:"path" `
	Data interface{} `json:"data" xml:"data" yaml:"data" form:"data" query:"data" `
}

func NewContext() *Context {
	return &Context{
		Path: []string{},
	}
}
