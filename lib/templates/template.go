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

	Buf *bytes.Buffer
}

func (T *Template) Render(data interface{}) ([]byte, error) {
	// endure we don't have nil, if so, there is a bug somewhere
	if T.T == nil {
		panic("template not set!")
	}

	var err error

	T.Buf.Reset()

	err = T.T.Execute(T.Buf, data)
	if err != nil {
		return nil, err
	}
	out := T.Buf.Bytes()

	return out, nil
}

// Creates a hof Template struct, initializing the correct template system. The system will be inferred if left empty
func CreateFromString(name, content string, delims *Delims) (t *Template, err error) {
	t = new(Template)
	t.Name = name
	t.Source = content

	// Golang wants helpers before parsing, and catches these errors early
	t.T = template.New(name)

	if delims != nil {
		t.T = t.T.Delims(delims.LHS, delims.RHS)
	}

	t.Buf = new(bytes.Buffer)

	t.AddGolangHelpers()

	t.T, err = t.T.Parse(content)

	return t, err
}
