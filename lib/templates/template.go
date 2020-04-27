package templates

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/aymerick/raymond"
)

type Template struct {
	// Original inputs
	Name   string
	Source string
	Config *Config

	// golang
	T *template.Template

	// mustache
	R *raymond.Template
}

func NewTemplate() *Template {
	return &Template{}
}

func NewMap() TemplateMap {
	return NewTemplateMap()
}

func (T *Template) Render(data interface{}) ([]byte, error) {
	// endure we don't have both ever, if so, there is a bug somewhere
	if T.T != nil && T.R != nil {
		panic("template instances are both set!")
	}

	// golang
	if T.T != nil {
		var b bytes.Buffer
		var err error

		err = T.T.Execute(&b, data)
		if err != nil {
			return nil, err
		}
		out := string(b.Bytes())

		out = T.Config.SwitchAfter(out)

		return []byte(out), nil
	}

	// mustache
	if T.R != nil {
		out, err := T.R.Exec(data)
		if err != nil {
			return nil, err
		}

		out = T.Config.SwitchAfter(out)

		return []byte(out), nil
	}

	return nil, fmt.Errorf("template instances are both empty")
}

// Creates a hof Template struct, initializing the correct template system. The system will be inferred if left empty
func CreateFromString(name, content, templateSystem string, config *Config) (t *Template, err error) {
	t = NewTemplate()
	t.Source = content
	t.Config = config

	if templateSystem == "" {
		templateSystem = inferTemplateSystem(content, t.Config)
	}

	switch templateSystem {
		case "golang":
			t.T, err = createGolangTemplate(name, content, t.Config)
			return t, err

		case "raymond":
			t.R, err = createRaymondTemplate(name, content, t.Config)
			return t, err

		default:
			return nil, fmt.Errorf("Unknown or unable to infer template system %q for file %q. Try setting explicitly", templateSystem, name)

	}

	return nil, err
}

func createGolangTemplate(name, content string, config *Config) (*template.Template, error) {

	// Golang wants helpers before parsing, and catches these errors early
	t := template.New(name)

	if config.LHS2_D != "{{" {
		t = t.Delims(config.LHS2_D, config.RHS2_D)
	}

	AddGolangHelpers(t)

	t, err := t.Parse(content)

	return t, err
}

func createRaymondTemplate(name, content string, config *Config) (*raymond.Template, error) {

	// Raymond want's to parse before helpers, and catches helper calls during exec
	content = config.SwitchBefore(content)

	r, err := raymond.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("While parsing file: %s\n%w\n", name, err)
	}

	AddRaymondHelpers(r)

	return r, nil
}

// TODO, we ought to be able to get config involved in this

func inferTemplateSystem(content string, config *Config) string {

	// hacky, but good enough for now, these ought to indicate raymond
	rayLhsCnt := strings.Count(content, config.LHS2_D + "#")
	rayLhsCnt += strings.Count(content, config.RHS2_D + "/")

	// pretty liberally assume golang
	if rayLhsCnt > 0 {
		fmt.Println("INFER RAYMOND")
		return "raymond"
	} else {
		fmt.Println("INFER GOLANG")
		return "golang"
	}

	// we shouldn't get here today, but maybe in the future?
	return "no template system detected"
}

