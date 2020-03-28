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
	F.SwitchDelimsBefore()

	fmt.Printf("template system: %q\n", sys)
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
	F.SwitchDelimsAfter()

	F.FormatRendered()

	return nil
}

func (F *File) SwitchDelimsTemplate(OLD, NEW string) {
	replace := F.TemplateContent

	replace = strings.ReplaceAll(replace, OLD, NEW)

	F.TemplateContent = replace
}

func (F *File) SwitchDelimsRendered(OLD, NEW string) {
	replace := F.RenderContent

	replace = bytes.ReplaceAll(replace, []byte(OLD), []byte(NEW))

	F.RenderContent = replace
}

func (F *File) SwitchDelimsBefore() {

	// Multi switch with temporary
	if F.AltDelims && F.SwapDelims {
		// Replace the swap or secondary with temp (this is the default for the template system)
		F.SwitchDelimsTemplate(F.LHS3_S, F.LHS3_T)
		F.SwitchDelimsTemplate(F.RHS3_S, F.RHS3_T)
		F.SwitchDelimsTemplate(F.LHS2_S, F.LHS2_T)
		F.SwitchDelimsTemplate(F.RHS2_S, F.RHS2_T)
	}

	if F.AltDelims {
		// Switch Swap for Default, which if you only set default, will work
		// do triple first, douvle second
		F.SwitchDelimsTemplate(F.LHS3_D, F.LHS3_S)
		F.SwitchDelimsTemplate(F.RHS3_D, F.RHS3_S)
		F.SwitchDelimsTemplate(F.LHS2_D, F.LHS2_S)
		F.SwitchDelimsTemplate(F.RHS2_D, F.RHS2_S)
	}

}

func (F *File) SwitchDelimsAfter() {

	// Multi switch undo
	if F.AltDelims && F.SwapDelims {
		// Undo the default to temp swap
		F.SwitchDelimsRendered(F.LHS3_T, F.LHS3_S)
		F.SwitchDelimsRendered(F.RHS3_T, F.RHS3_S)
		F.SwitchDelimsRendered(F.LHS2_T, F.LHS2_S)
		F.SwitchDelimsRendered(F.RHS2_T, F.RHS2_S)

	}

	// shouldn't have to do anything since we rendered by now
	//   and there should have been only one template system
	//   or  another that is not effected by renering

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
