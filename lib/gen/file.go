package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/epiclabs-io/diff3"
	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/hofstadter-io/hof/lib/templates"
)

type File struct {
	// Input Data, local to this file
	In           map[string]interface{}

  // The full path under the output location
  // empty implies don't generate, even though it may endup in the list
	Filepath     string

	// Template parameters
	TemplateSystem string  // which system ['text/template'(default), 'mustache']
	Template       string  // The content, takes precedence over next option
	TemplateName   string  // Named template

  //
  // Template delimiters
  //
  //   these are for advanced usage, you shouldn't have to modify them normally

  // Alt and Swap Delims,
	//   becuase the defaulttemplate systems use `{{` and `}}`
	//   and you may choose to use other delimiters, but the lookup system is still based on the template system
  //   and if you want to preserve those, we need three sets of delimiters
  AltDelims  bool
  SwapDelims bool

  // The default delimiters
  // You should change these when using alternative style like jinjas {% ... %}
  // They also need to be different when using the swap system
  LHS2_D string
  RHS2_D string
  LHS3_D string
  RHS3_D string

  // These are the same as the default becuase
  // the current template systems require these.
  //   So these should really never change or be overriden until there is a new template system
  //     supporting setting the delimiters dynamicalldelimiters dynamicallyy
  LHS2_S string
  RHS2_S string
  LHS3_S string
  RHS3_S string

  // The temporary delims to replace swap with while also swapping
  // the defaults you set to the swap that is required by the current templet systems
	// You need this when you are double templating a file and the top-level system is not the default
  LHS2_T string
  RHS2_T string
  LHS3_T string
  RHS3_T string

	//
	// Hof internal usage
	//

	// Content
	TemplateContent string
	RenderContent   []byte
	FinalContent    []byte

	// Shadow related
	ShadowFile *File
	UserFile   *File

	DoWrite bool

	// Bookkeeping
	Errors []error
	FileStats
}

func (F *File) Render() error {
	var err error

	// TODO eventually look for template file by file name
	// in some cache, but do this somewhere else, so that
	// we have an abstract template system
	if F.TemplateContent == "" {
		F.TemplateContent = F.Template
	}

	err = F.RenderTemplate()
	if err != nil {
		return err
	}

	// Check to see if they are the same, if so, then "skip"
	// fmt.Println(F.Filepath, len(F.RenderContent), F.ShadowFile)
	if F.ShadowFile != nil {
		F.ReadShadow()
		if bytes.Compare(F.RenderContent, F.ShadowFile.FinalContent) == 0 {
			// Let's check if there is a user file or not
			_, err := os.Lstat(F.Filepath)
			if err != nil {
				// make sure we check err for something actually bad
				if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" {
					return err
				}
				// file does not exist
				F.IsNew = 1
				F.DoWrite = true
				F.FinalContent = F.RenderContent
				return nil
			}
			F.IsSame = 1
			return nil
		}
	}

	// Possibly read user
	F.ReadUser()

	// figure out if / how to merge and produce final content
	F.DoWrite, err = F.UnifyContent()
	if err != nil {
		F.IsErr = 1
		return err
	}

	return nil
}

func (F *File) ReadUser() error {

	_, err := os.Lstat(F.Filepath)
	if err != nil {
		// make sure we check err for something actually bad
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" {
			return err
		}
		return nil
	}

	content, err := ioutil.ReadFile(F.Filepath)
	if err != nil {
		return err
	}

	F.UserFile = &File {
		Filepath: F.Filepath,
		FinalContent: content,
	}

	return nil
}

func (F *File) UnifyContent() (write bool, err error) {
	// set this first, possible change later in this function
	F.FinalContent = F.RenderContent

	// If there is a user file...
	if F.UserFile != nil {
		if F.ShadowFile != nil {
			// Need to compare all 3
			// But first a shortcut
			if bytes.Compare(F.UserFile.FinalContent, F.ShadowFile.FinalContent) == 0 {
				// fmt.Println("User == Shadow", len(F.UserFile.FinalContent), len(F.ShadowFile.FinalContent))
				// Just write it out, no user modifications
				F.IsModified = 1
				F.IsModifiedRender = 1
				return true, nil
			}

			O := bytes.NewReader(F.ShadowFile.FinalContent)
			A := bytes.NewReader(F.UserFile.FinalContent)
			B := bytes.NewReader(F.FinalContent)
			labelA := "Your File"
			labelB := "New File"
			detailed := true

			result, err := diff3.Merge(A, O, B, detailed, labelA, labelB)
			if err != nil {
				F.IsErr = 1
				return false, err
			}

			merged, err := ioutil.ReadAll(result.Result)
			if err != nil {
				F.IsErr = 1
				return false, err
			}

			if result.Conflicts {
				F.IsConflicted = 1
			}

			F.IsModified = 1
			F.IsModifiedDiff3 = 1
			F.FinalContent = merged

			return true, nil

		} else {

			// Compare new content to User content
			if bytes.Compare(F.RenderContent, F.UserFile.FinalContent) == 0 {
				// Don't write it out, no user modifications, or the same modifications?
				F.IsSame = 1
				return false, nil

			} else {
				// 2-way diff, the user made modifications
				dmp := diffmatchpatch.New()
				// Do this backwards, how do we get from user file to the new one
				diffs := dmp.DiffMain(string(F.FinalContent), string(F.UserFile.FinalContent), false)

				// Now skip anything the user "deleted" from the file, i.e. new content
				for _, d := range diffs {
					if d.Type == -1 {
						// "skip" by setting equal, otherwise we mess things up by not including it
						d.Type = 0
					}
				}

				merged := dmp.DiffText2(diffs)
				F.IsModified = 1
				F.IsModifiedOutput = 1
				F.FinalContent = []byte(merged)

				return true, nil
			}
		}
	}

	// Otherwise, this is a new file
	F.IsNew = 1

	return true, nil
}

func (F *File) RenderTemplate() error {

	// Will check to see what the situation is
	F.SwitchDelimsBefore()

	sys := strings.ToLower(F.TemplateSystem)
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
