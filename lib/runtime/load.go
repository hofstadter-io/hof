package runtime

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/cue/token"
	"cuelang.org/go/encoding/json"
	"cuelang.org/go/encoding/yaml"

	// "github.com/kr/pretty"

	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/hof"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func (R *Runtime) Load() (err error) {
	if R.Flags.Verbosity > 0 {
		fmt.Println("Loading from:", R.Entrypoints)
	}

	// stats / timing
	start := time.Now()
	defer func() {
		end := time.Now()
		R.Stats.Add("cue/load", end.Sub(start))
	}()

	R.prepPlacedDatafiles()
	
	err = R.load()
	if err != nil {
		return err
	}

	if !R.Flags.IngoreErrors {
		err = R.Value.Validate()
		if err != nil {
			return cuetils.ExpandCueError(err)
		}
	}

	R.Nodes, err = hof.FindHofsOrig(R.Value)
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) prepPlacedDatafiles() {
	start := time.Now()
	defer func() {
		end := time.Now()
		R.Stats.Add("data/load", end.Sub(start))
	}()
	R.origEntrypoints = make([]string, 0, len(R.Entrypoints))
	entries := []string{}

	// we cloud probably do something to support both globs and @ placement, but there is -l as well
	for _, E := range R.Entrypoints {
		R.origEntrypoints = append(R.origEntrypoints, E)

		// expand globs
		if strings.Contains(E, "*") {
			files, err := yagu.FilesFromGlobs([]string{E})
			if err != nil {
				fmt.Println("warning: error while globing %q: %v", E, err)
			}
			entries = append(entries, files...)
			continue
		}

		// placed data
		if strings.Contains(E, "@") {
			parts := strings.Split(E, "@")
			if len(parts) == 2 {
				// add the mapping
				fname, fpath := parts[0], parts[1]
				R.dataMappings[fname] = fpath
				entries = append(entries, fname)
			continue
			}
		}

		// else, just the entrypoint as is
		entries = append(entries, E)
	}

	if R.Flags.Verbosity > 1 {
		fmt.Println("update entrypoints:", entries)
	}

	R.Entrypoints = entries
}

func (R *Runtime) load() (err error) {
	beg := time.Now()
	defer func() {
		end := time.Now()
		R.Stats.Add("gen/load", end.Sub(beg))
	}()

	// XXX TODO XXX
	//  add the second arg from our runtime when implemented?
	//  is this to support multiple R's at oncce?
	//  or do we just wait for CUE to be better?
	if R.CueContext == nil {
		R.CueContext = cuecontext.New()
	}
	// fmt.Printf("%# v\n", pretty.Formatter(R.CueConfig))
	R.CueConfig.DataFiles = R.Flags.IncludeData
	R.BuildInstances = load.Instances(R.Entrypoints, R.CueConfig)

	if l := len(R.BuildInstances); l == 0 {
		return fmt.Errorf("expected at least one build instance, got none", l)
	} else if l >= 2 {
		// this looks to always be empty when it is created, so we just ignore it
		// fmt.Printf("warning, go more than one instance: %#v %#v\n", R.BuildInstances[0], R.BuildInstances[1])
	}

	// we always take the first build instance
	bi := R.BuildInstances[0]

	if bi.Err != nil {
		return bi.Err
	}
	if bi.Incomplete {
		return fmt.Errorf("incomplete build instance, ask devs")
	}

	err = R.prepOrphanedFiles(bi)
	if err != nil {
		return err
	}

	// Build the Instance
	R.Value = R.CueContext.BuildInstance(bi)

	// unify any -I inputs
	for i, I := range R.Flags.InputData {
		if strings.Contains(I, "=") {
			parts := strings.Split(I, "=")
			R.Value = R.Value.FillPath(cue.ParsePath(parts[0]),parts[1])
		} else {
			v := R.CueContext.CompileString(I)
			if v.Err() != nil {
				err := cuetils.ExpandCueError(v.Err())
				return fmt.Errorf("in -I(%d) flag '%s': %w", i, I, err)
			}
			R.Value = R.Value.FillPath(cue.ParsePath(""), v)
		}
	}

	return nil
}

func (R *Runtime) prepOrphanedFiles(bi *build.Instance) (err error) {
	// a bit hacky...
	if R.DontPlaceOrphanedFiles {
		return nil
	}
	// TODO, compare with len(entrypoints)
	//       and be more intelligent
	var errs []errors.Error

	if R.Flags.Verbosity > 1 {
		fmt.Println("ID:", bi.ID(), bi.PkgName, bi.Module)
	}
	pkg := bi.PkgName
	//if bi.Module == "" {
	//  pkg = bi.ID()
	//}


	// handle data files
	for i, f := range bi.OrphanedFiles {
		// this function also checks to see if we should include the file
		//   based on a few settings, but we have to do some path handling first...
		F, err := R.LoadOrphanedFile(f, pkg, bi.Root, bi.Dir, i, len(bi.OrphanedFiles))
		if err != nil {
			if R.Flags.Verbosity > 1 {
				fmt.Println("[load] error in data:", f.Filename, err)
			}
			errs = append(errs, errors.Promote(err,""))
			continue
		}
		// we don't know what this file is
		if F == nil {
			if R.Flags.Verbosity > 1 {
				fmt.Println("[load] ignoring data:", f.Filename)
			}
			continue
		}
		if R.Flags.Verbosity > 1 {
			fmt.Println("[load] including data:", f.Filename)
		}

		// embed the data file, already placed if needed
		bi.AddSyntax(F)
	}

	if len(errs) > 0 {
		_e := errors.New("in prepOrphanedFiles")
		e := errors.Promote(_e,"")
		for _, err := range errs {
			e = errors.Append(e, err)
		}
		return e
	}

	return nil
}

func (R *Runtime) LoadOrphanedFile(f *build.File, pkgName string, root, dir string, index, total int) (F *ast.File, err error) {
	if R.Flags.Verbosity > 1 {
		fmt.Println("[load]:", f.Filename, reflect.TypeOf(f.Source))
	}

	var d []byte

	fname := f.Filename
	// strip dir from fname
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	fname = strings.TrimPrefix(fname, dir)


	// only load data files which are explicitly listed
	//   or if the --include-data flag is set
	// this is checking to see if we should return early
	if !R.Flags.IncludeData { // user is not including all data
		// so check if explicitly supplied as an arg
		// fmt.Println("checking:", fname, R.Entrypoints)
		match := false
		for _, e := range R.Entrypoints {
			if filepath.Clean(fname) == filepath.Clean(e) {
				match = true
				break
			}
		}
		if !match {
			return nil, nil
		}
	}

	mapping := R.dataMappings[fname]

	if R.Flags.Verbosity > 1 && mapping != "" {
		fmt.Printf("[load] found entrypoint mapping: %s -> %s\n", f.Filename, mapping)
	}

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

	withContext := func(e ast.Expr) ast.Expr {
		if R.Flags.WithContext {
			return ast.NewStruct(
				"data", e,
				"filename", ast.NewString(f.Filename),
				"index", ast.NewLit(token.INT, strconv.Itoa(index)),
				"recordCount", ast.NewLit(token.INT, strconv.Itoa(total)),
			)
		}
		return e
	}


	switch f.Encoding {

	case "json":
		A, err := json.Extract(f.Filename, d)
		if err != nil {
			return nil, fmt.Errorf("while extracting json file: %w", err)
		}

		C := withContext(A)

		A, err = R.placeOrphanInAST(A.(*ast.StructLit), C, mapping)
		if err != nil {
			return nil, err
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

		A := ast.NewStruct()
		A.Elts = F.Decls
		C := withContext(A)

		A, err = R.placeOrphanInAST(A, C, mapping)
		if err != nil {
			return nil, err
		}

		F.Decls = []ast.Decl{A}

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
		// otherwise we should ignore other files implicitly discovered
		// we see go files here for example

		// todo, re-enable this with better checks
		// err := fmt.Errorf("unknown encoding for %q %q", f.Filename, f.Encoding)
		// return nil, err
		return nil, nil
	}

}

func (R *Runtime) placeOrphanInAST(N ast.Node, C ast.Expr, mapping string) (*ast.StructLit, error) {
	S := N.(*ast.StructLit)
	// fmt.Println("GOT HERE", S, R.Flags.Path)
	if mapping != "" {
		// @path placed datafiles
		ps := cue.ParsePath(mapping).Selectors()
		// go in reverse, so we build up a tree
		for i := len(ps)-1; i >= 0; i--  {
			// build our label from the mapping path
			p := ps[i]
			ident := ast.NewIdent(p.String())

			// create a struct with a field
			f := &ast.Field {
				Label: ident,
				Value: S,
			}
			s := ast.NewStruct(f)

			// now update
			S = s
		}
	} else if len(R.Flags.Path) > 0 {
		// -l/--path placed datafiles
		ps := R.Flags.Path
		// fmt.Println("PathFlags:", ps)

		ctx := R.CueContext
		v := ctx.BuildExpr(C)
		if v.Err() != nil {
			return nil, v.Err()
		}

		for i := len(ps)-1; i >= 0; i--  {
			// build our label from the mapping path
			p := ps[i]

			pv := ctx.CompileString(
				p, 
				cue.Filename(p),
				cue.InferBuiltins(true),
				cue.Scope(v),
			)
			if pv.Err() != nil {
				return nil, pv.Err()
			}

			str, err := pv.String()
			if err != nil {
				return nil, err
			}

			ident := ast.NewIdent(str)

			// create a struct with a field
			f := &ast.Field {
				Label: ident,
				Value: S,
			}
			s := ast.NewStruct(f)

			// now update
			S = s
		}
	}

	return S, nil

}
