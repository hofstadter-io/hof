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

	// when CUE entrypoints have @placement
	origEntrypoints []string

	// when a user supplies an data.json@path.to.field
	dataMappings    map[string]string
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
	CRT.prepPlacedDatafiles()
	return CRT.load()
}

func (CRT *CueRuntime) prepPlacedDatafiles() {
	CRT.origEntrypoints = make([]string, 0, len(CRT.Entrypoints))

	for i, E := range CRT.Entrypoints {
		CRT.origEntrypoints = append(CRT.origEntrypoints, E)
		if !strings.Contains(E, "@") {
			continue
		}

		parts := strings.Split(E, "@")
		if len(parts) != 2 {
			continue
		}

		// add the mapping
		fname, fpath := parts[0], parts[1]
		CRT.dataMappings[fname] = fpath

		CRT.Entrypoints[i] = fname
	}

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
	CRT.BuildInstances = load.Instances(CRT.Entrypoints, CRT.CueConfig)
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
			F, err := CRT.loadOrphanedFile(f, bi.PkgName, bi.Root, bi.Dir)
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

func (CRT *CueRuntime) loadOrphanedFile(f *build.File, pkgName string, root, dir string) (F *ast.File, err error) {

	var d []byte

	fname := f.Filename
	// strip dir from fname
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	fname = strings.TrimPrefix(fname, dir)

	mapping := CRT.dataMappings[fname]
	// if mapping != "" {
	// 	fmt.Printf("found entrypoint mapping: %s -> %s\n", f.Filename, mapping)
	// }

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

		if mapping != "" {
			ps := cue.ParsePath(mapping).Selectors()
			// go in reverse, so we build up a tree
			for i := len(ps)-1; i >= 0; i--  {
				// build our label from the mapping path
				p := ps[i]
				ident := ast.NewIdent(p.String())

				// create a struct with a field
				f := &ast.Field {
					Label: ident,
					Value: A,
				}
				s := ast.NewStruct(f)

				// now update
				A = s
			}
		}

		// add a package decl so the data is referencable from the cue
		pkgDecl := &ast.Package {
			Name: ast.NewIdent(pkgName),
		}

		// extract the json top level fields (removing the outer unnamed struct)
		jsonDecls := []ast.Decl{pkgDecl, A}
		if mapping == "" {
			switch a := A.(type) {
				case *ast.StructLit:
					jsonDecls = append([]ast.Decl{pkgDecl}, a.Elts...)
			}
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

		if mapping != "" {
			A := ast.NewStruct()
			A.Elts = F.Decls
			ps := cue.ParsePath(mapping).Selectors()
			// go in reverse, so we build up a tree
			for i := len(ps)-1; i >= 0; i--  {
				// build our label from the mapping path
				p := ps[i]
				ident := ast.NewIdent(p.String())

				// create a struct with a field
				f := &ast.Field {
					Label: ident,
					Value: A,
				}
				s := ast.NewStruct(f)

				// now update
				A = s
			}

			F.Decls = []ast.Decl{A}
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
