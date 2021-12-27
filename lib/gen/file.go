package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"strings"

	"github.com/epiclabs-io/diff3"
	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/hofstadter-io/hof/lib/templates"
)

type File struct {
	// Input Data, local to this file
	In map[string]interface{}

	// The full path under the output location
	// empty implies don't generate, even though it may endup in the list
	Filepath string

	// Template parameters
	TemplateContent string // The content, takes precedence over next option
	TemplatePath    string // Named template
	CueFile         bool
	Concrete        bool

	// Template delimiters
	TemplateDelims *templates.Delims

	//
	// Hof internal usage
	//

	// Generator that owns this file
	// TODO, does this break with the multi-generator hack?
	//       or is this the behavior defined by the way generators are included together
	//       because the user will be able to do it both ways, the hack doesn't to away,
	// there's just a more recommended way to do it
	Gen *Generator

	// Template Instance Pointer
	//   If local, this will be created when the template content is laoded
	//   If a named template, acutal template lives in the generator and is created at folder import time
	TemplateInstance *templates.Template

	// Content
	RenderContent []byte
	FinalContent  []byte

	// Shadow related
	ShadowFile *File
	UserFile   *File

	DoWrite bool

	// Bookkeeping
	Errors []error
	FileStats
}

func (F *File) Render(shadow_basedir string) error {
	var err error

	err = F.RenderTemplate()
	if err != nil {
		F.FileStats.IsErr = 1
		F.Errors = append(F.Errors, err)
		return err
	}
	// fmt.Println("   rendered:", F.Filepath, len(F.RenderContent))

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

	F.UserFile = &File{
		Filepath:     F.Filepath,
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
	var err error

	F.RenderContent, err = F.TemplateInstance.Render(F.In)
	if err != nil {
		return err
	}

	err = F.FormatRendered()
	if err != nil {
		fmt.Printf("---- Rendering error for template: %q output: %q content:\n", F.TemplatePath, F.Filepath)
		fmt.Println(string(F.RenderContent))
		fmt.Println("----")
		return err
	}

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
	}

	return nil
}
