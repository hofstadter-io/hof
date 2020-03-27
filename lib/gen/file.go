package gen

import (
	"bytes"
	// "fmt"
	"io/ioutil"
	"os"

	"github.com/epiclabs-io/diff3"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type File struct {
	// Inputs
	Filename string
	Template string
	In       map[string]interface{}

	// Template parameters
	TemplateSystem string
	// Delimiter Values
	// ...TODO

	// Content
	TemplateContent string
	RenderContent   []byte
	FinalContent    []byte

	// Shadow related
	ShadowFile *File
	UserFile   *File

	// Bookkeeping
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
	// fmt.Println(F.Filename, len(F.RenderContent), F.ShadowFile)
	if F.ShadowFile != nil {
		F.ReadShadow()
		if bytes.Compare(F.RenderContent, F.ShadowFile.FinalContent) == 0 {
			F.IsSame = 1
			return nil
		}
	}

	// Possibly read user
	F.ReadUser()

	// figure out if / how to merge
	doWrite, err := F.UnifyContent()
	if err != nil {
		F.IsErr = 1
		return err
	}

	if doWrite {
		F.WriteOutput()
		F.WriteShadow()
	}

	return nil
}

func (F *File) ReadUser() error {

	_, err := os.Lstat(F.Filename)
	if err != nil {
		// make sure we check err for something actually bad
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" {
			return err
		}
		return nil
	}

	content, err := ioutil.ReadFile(F.Filename)
	if err != nil {
		return err
	}

	F.UserFile = &File {
		Filename: F.Filename,
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

