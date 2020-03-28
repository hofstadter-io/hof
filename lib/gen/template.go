package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/hofstadter-io/hof/lib/templates"
)

func (F *File) RenderTemplate() error {
	sys := strings.ToLower(F.TemplateSystem)

	// Will check to see what the situation is
	F.SwapDelimsBefore()

	switch sys {
		case "mustache", "raymond", "handlebars":
			F.RenderRaymondTemplate()

		// sudo default
		case "", "golang", "text", "text/template":
			F.RenderGolangTemplate()

		default:
			return fmt.Errorf("Unknown template system: ", sys)
	}

	// Will check to see what the situation is
	F.SwapDelimsAfter()

	F.FormatRendered()

	return nil
}

func (F *File) SwapDelimsBefore() error {

	return nil
}

func (F *File) SwapDelimsAfter() error {

	return nil
}

func (F *File) RenderGolangTemplate() error {

	t := template.Must(template.New(F.Filepath).Parse(F.TemplateContent))

	var b bytes.Buffer
	var err error

	err = t.Execute(&b, F.In)
	if err != nil {
		return err
	}

	F.RenderContent = b.Bytes()

	return nil
}

func (F *File) RenderRaymondTemplate() error {

	t, err := templates.CreateTemplateFromString(F.Filepath, F.TemplateContent)
	if err != nil {
		return err
	}

	out, err := t.Render(F.In)
	if err != nil {
		return err
	}

	F.RenderContent = []byte(out)

	return nil
}


func (F *File) FormatRendered() error {

	// If Golang only
	if strings.HasSuffix(F.Filepath, ".go") {
		fmtd, err := format.Source(F.RenderContent)
		if err != nil {
			return err
		}

		F.RenderContent = fmtd
	} else {
	}

	return nil
}
