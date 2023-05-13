package gen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/hofstadter-io/hof/lib/diff3"
	"github.com/hofstadter-io/hof/lib/templates"
)

type File struct {
	// Input Data, local to this file
	// this can sometimes be a list for Adhoc
	In any

	// The full path under the output location
	// empty implies don't generate, even though it may endup in the list
	Filepath string

	// Template parameters (only one should be set at a time i.e. != "")
	TemplateContent string // The content, takes precedence over next option
	TemplatePath    string // Named template
	DatafileFormat  string // Data format file
	StaticFile      bool

	// Formatting
	FormattingDisabled bool
	FormattingConfig   *FmtConfig

	// Template delimiters
	TemplateDelims *templates.Delims

	// CUE settings
	Package        string
	Raw            bool
	Final          bool
	Concrete       bool
	Definitions    bool
	Optional       bool 
	Hidden         bool 
	Attributes     bool
	Docs           bool 
	InlineImports  bool 
	ErrorsAsValues bool



	//
	// Hof internal usage
	//

	// CUE value for datafiles
	// (we use a different name here so that it does not automatically try to decode, which would require concreteness)
	Value cue.Value

	// Template Instance Pointer
	//   If local, this will be created when the template content is loaded
	//   If a named template, actual template lives in the generator and is created at folder import time
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
	parent *Generator
}

func (F *File) MergeIn(other map[string]any) {
	// merge with other if map / struct
	// (this is probably always the case (should be enforced?))
	switch fin := F.In.(type) {
	case map[string]any:
		for k, v := range other {
			// only copy in top-level elements which do not exist already
			if _, ok := fin[k]; !ok {
				fin[k] = v
			}
		}
		F.In = fin
	}
}

func (F *File) Render(outdir string, UseDiff3, NoFmt bool) error {
	// fmt.Println("F.Render:", F.Filepath, UseDiff3)
	var err error
	F.RenderContent = bytes.TrimSpace(F.RenderContent)

	if F.DatafileFormat != "" {
		err = F.RenderData()
		if err != nil {
			err = fmt.Errorf("In: %q %w", F.Filepath, err)
			F.FileStats.IsErr = 1
			F.Errors = append(F.Errors, err)
			return err
		}
	} else if F.StaticFile {

	} else {
		err = F.RenderTemplate(NoFmt)
		if err != nil {
			F.FileStats.IsErr = 1
			F.Errors = append(F.Errors, err)
			return err
		}
	}
	// fmt.Println("   rendered:", F.Filepath, len(F.RenderContent), F.ShadowFile != nil)

	// Check to see if they are the same, if so, then "skip"
	if UseDiff3 && F.ShadowFile != nil {
		if bytes.Compare(F.RenderContent, F.ShadowFile.FinalContent) == 0 {
			// Let's check if there is a user file or not
			_, err := os.Lstat(filepath.Join(outdir, F.Filepath))
			if err != nil {
				// make sure we check err for something actually bad
				if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" {
					return err
				}
				// fmt.Println("  new file")
				// file does not exist
				F.IsNew = 1
				F.DoWrite = true
				F.FinalContent = F.RenderContent
				return nil
			}
			// fmt.Println("  same file")
			F.DoWrite = false
			F.IsSame = 1
			return nil
		}
	}

	// Possibly read user
	if UseDiff3 {
		// fmt.Println("  read user")
		F.ReadUser(outdir)
	}

	// figure out if / how to merge and produce final content
	F.DoWrite, err = F.UnifyContent(UseDiff3)
	if err != nil {
		F.IsErr = 1
		return err
	}

	return nil
}

// read the file contents relative to the output dir
func (F *File) ReadUser(outdir string) error {
	fp := filepath.Join(outdir, F.Filepath)

	content, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	content = bytes.TrimSpace(content)

	F.UserFile = &File{
		Filepath:     fp,
		FinalContent: content,
	}

	return nil
}

func (F *File) UnifyContent(UseDiff3 bool) (write bool, err error) {
	// fmt.Println("unify:", F.Filepath)
	// set this first, possible change later in this function
	F.FinalContent = bytes.TrimSpace(F.RenderContent)

	// If there is a user file...
	if UseDiff3 && F.UserFile != nil {
		// must have all 3
		if F.ShadowFile != nil {
			return F.diff3()
		}
		return F.diff2()
	} // end UseDiff3

	// Otherwise, this is a new file
	F.IsNew = 1

	return true, nil
}

func (F *File) diff3() (write bool, err error) {
	// fmt.Println("diff3:", F.Filepath)

	FC := F.FinalContent
	UF := F.UserFile.FinalContent
	SF := F.ShadowFile.FinalContent

	// if shadow and user content same
	// Just write it out, no user modifications ever
	if bytes.Compare(UF, SF) == 0 {
		F.IsModified = 1
		F.IsModifiedRender = 1
		return true, nil
	}

	//merged := diff3.Merge(string(SF), string(UF), string (FC))
	//has1 := strings.Contains(merged,diff3.Sep1)
	//has2 := strings.Contains(merged,diff3.Sep2)
	//has3 := strings.Contains(merged,diff3.Sep3)
	//if has1 && has2 && has3 {
		//F.IsConflicted = 1
	//}
	//merged = strings.TrimSpace(merged)

	// Now need to compare all 3
	labelA := "Your File"
	A := bytes.NewReader(UF)
	O := bytes.NewReader(SF)
	B := bytes.NewReader(FC)
	labelB := "Code Gen"
	detailed := false

	result, err := diff3.Merge(A, O, B, detailed, labelA, labelB)
	if err != nil {
		F.IsErr = 1
		return false, err
	}

	// TODO, filter results which have an empty "Your File" content
	// (will this break other legit cases?)

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
	F.FinalContent = []byte(merged)

	return true, nil
}

func (F *File) diff2() (write bool, err error) {
	// probably the first time we gen with diff enable
	// fmt.Println("diff2:", F.Filepath)
	// fmt.Println("GOT HERE, tell devs")

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
		F.FinalContent = bytes.TrimSpace([]byte(merged))

		return true, nil
	}
}

func (F *File) RenderData() (err error) {
	F.RenderContent, err = F.formatData(F.Value, F.DatafileFormat)
	if err != nil {
		F.Errors = append(F.Errors, err)
		return err
	}

	return nil
}

func (F *File) RenderTemplate(nofmt bool) (err error) {

	F.RenderContent, err = F.TemplateInstance.Render(F.In)
	if err != nil {
		F.Errors = append(F.Errors, err)
		return err
	}

	if nofmt {
		return
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
