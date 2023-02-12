package cuetils

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/json"
	"cuelang.org/go/encoding/yaml"
)

type CueSyntaxOptions struct {
	Attributes  bool
	Concrete    bool
	Definitions bool
	Docs        bool
	Hidden      bool
	Optional    bool
}

func (CSO CueSyntaxOptions) MakeOpts() []cue.Option {
	return []cue.Option{
		cue.Attributes(CSO.Attributes),
		cue.Concrete(CSO.Concrete),
		cue.Definitions(CSO.Definitions),
		cue.Docs(CSO.Docs),
		cue.Hidden(CSO.Hidden),
		cue.Optional(CSO.Optional),
	}
}

var (
	DefaultSyntaxOpts = CueSyntaxOptions{
		Attributes:  true,
		Concrete:    false,
		Definitions: true,
		Docs:        true,
		Hidden:      true,
		Optional:    true,
	}
)

type CueRuntime struct {
	Entrypoints []string
	Workspace   string
	FS          billy.Filesystem

	CueContext     *cue.Context
	CueConfig      *load.Config
	BuildInstances []*build.Instance
	CueErrors      []error
	FieldOpts      []cue.Option

	CueInstance *cue.Instance
	CueValue    cue.Value
	Value       interface{}
}

func (CRT *CueRuntime) ConvertToValue(in interface{}) (cue.Value, error) {
	O, ook := in.(cue.Value)
	if !ook {
		switch T := in.(type) {
		case string:
			i := CRT.CueContext.CompileString(T)
			if i.Err() != nil {
				return O, i.Err()
			}
			v := i.Value()
			if v.Err() != nil {
				return v, v.Err()
			}
			O = v

		default:
			return O, fmt.Errorf("unknown type %v in convertToValue(in)", T)
		}
	}

	return O, nil
}

func (CRT *CueRuntime) Load() (err error) {
	return CRT.load()
}

func (CRT *CueRuntime) load() (err error) {
	// possibly, check for workpath
	if CRT.Workspace != "" {
		_, err = os.Lstat(CRT.Workspace)
		if err != nil {
			if _, ok := err.(*os.PathError); !ok && (strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file")) {
				// error is worse than non-existant
				return err
			}
			// otherwise, does not exist, so we should init?
			// XXX want to let applications decide how to handle this
			return err
		}
	}

	var errs []error

	// XXX TODO XXX
	//  add the second arg from our runtime when implemented
	if CRT.CueContext == nil {
		CRT.CueContext = cuecontext.New()
	}
	CRT.BuildInstances = load.Instances(CRT.Entrypoints, nil)
	for _, bi := range CRT.BuildInstances {
		// fmt.Printf("%d: start\n", i)

		if bi.Err != nil || bi.Incomplete {
			es := errors.Errors(bi.Err)
			for _, e := range es {
				errs = append(errs, e.(error))
			}
			continue
		}

		// TODO, compare with len(entrypoints)
		//       and be more intelligent

		// handle data files
		for _, f := range bi.OrphanedFiles {
			F, err := loadOrphanedFile(f, bi.PkgName)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			bi.AddSyntax(F)
		}

		// Build the Instance
		V := CRT.CueContext.BuildInstance(bi)
		if V.Err() != nil {
			es := errors.Errors(V.Err())
			for _, e := range es {
				errs = append(errs, e.(error))
			}
			continue
		}

		CRT.CueValue = V

	}

	if len(errs) > 0 {
		CRT.CueErrors = errs
		s := fmt.Sprintf("Errors while loading Cue entrypoints: %s %v\n", CRT.Workspace, CRT.Entrypoints)
		for _, E := range errs {
			es := errors.Errors(E)
			for _, e := range es {
				s += CueErrorToString(e)
			}
		}
		return fmt.Errorf(s)
	}

	return nil
}

func loadOrphanedFile(f *build.File, pkgName string) (F *ast.File, err error) {

	var d []byte

	if f.Filename == "-" {
    reader := bufio.NewReader(os.Stdin)
    var buf bytes.Buffer
    for {
        b, err := reader.ReadByte()
        if err != nil {
            break
        }
        buf.WriteByte(b)
    }
		d = buf.Bytes()
	} else {
		d, err = os.ReadFile(f.Filename)
		if err != nil {
			return nil, fmt.Errorf("while loading data file: %w", err)
		}
	}


	switch f.Encoding {

	case "json":
		A, err := json.Extract(f.Filename, d)
		if err != nil {
			return nil, fmt.Errorf("while extracting json file: %w", err)
		}

		// add a package decl so the data is referencable from the cue
		pkgDecl := &ast.Package {
			Name: ast.NewIdent(pkgName),
		}

		// extract the json top level fields (removing the outer unnamed struct)
		jsonDecls := []ast.Decl{pkgDecl, A}
		switch a := A.(type) {
			case *ast.StructLit:
				jsonDecls = append([]ast.Decl{pkgDecl}, a.Elts...)
		}

		// construct an ast.File to be consistent with Yaml
		// (and also provide a filename for errors)
		F := &ast.File{
			Filename: f.Filename,
			// Decls: A.(*ast.StructLit).Elts,
			// Decls:    []ast.Decl{pkgDecl, A},
			Decls: jsonDecls,
		}

		return F, nil

	case "yml", "yaml":
		F, err := yaml.Extract(f.Filename, d)
		if err != nil {
			return nil, fmt.Errorf("while extracting yaml file: %w", err)
		}

		// add a package decl so the data is referencable from the cue
		pkgDecl := &ast.Package {
			Name: ast.NewIdent(pkgName),
		}
		F.Decls = append([]ast.Decl{pkgDecl}, F.Decls...)

		// merge in data
		return F, nil

	// ....
	// case: ...

	// TODO, handle other formats (toml,json)
	//       which hof already works with

	default:
		// should only do this if it was also an arg
		// otherwise we should ignore other files implicity discovered

		// todo, re-enable this with better checks
		// err := fmt.Errorf("unknown encoding for %q %q", f.Filename, f.Encoding)
		// return nil, err
		return nil, nil
	}

}
