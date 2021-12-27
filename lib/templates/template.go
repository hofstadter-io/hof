package templates

import (
	"bytes"
	"text/template"
)

type Delims struct {
	LHS string
	RHS string
}

type Template struct {
	// Original inputs
	Name   string
	Source string
	Delims *Delims

	// golang
	T *template.Template
}

func NewTemplate() *Template {
	return &Template{}
}

func (T *Template) Render(data interface{}) ([]byte, error) {
	// endure we don't have nil, if so, there is a bug somewhere
	if T.T == nil {
		panic("template not set!")
	}

	var b bytes.Buffer
	var err error

	err = T.T.Execute(&b, data)
	if err != nil {
		return nil, err
	}
	out := string(b.Bytes())

	return []byte(out), nil
}

// Creates a hof Template struct, initializing the correct template system. The system will be inferred if left empty
func CreateFromString(name, content string, delims *Delims) (t *Template, err error) {
	t = NewTemplate()
	t.Name = name
	t.Source = content

	// Golang wants helpers before parsing, and catches these errors early
	t.T = template.New(name)

	if delims != nil {
		t.T = t.T.Delims(delims.LHS, delims.RHS)
	}

	AddGolangHelpers(t.T)

	t.T, err = t.T.Parse(content)

	return t, err
}
