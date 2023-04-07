package cmd

import (
	"os"

	/*
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/ast/astutil"
	"cuelang.org/go/cue/format"
	*/

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func info(R *runtime.Runtime, dflags flags.DatamodelPflagpole) error {

	// find max label width after indentation for column alignment
	max := findMaxLabelLen(R, dflags)

	for _, dm := range R.Datamodels {
		if err := dm.PrintInfo(os.Stdout, max, dflags); err != nil {
			return err
		}

		/*

		name := dm.Hof.Label
		p := cue.ParsePath(name)

		ctx := dm.Value.Context()
		val := ctx.CompileString("_")

		val = val.FillPath(p, dm.Value)

		// add lacunas

		node := val.Syntax(
			cue.Final(),
			cue.Docs(true),
			cue.Attributes(true),
			cue.Definitions(true),
			cue.Optional(true),
			cue.Hidden(true),
			cue.Concrete(true),
			cue.ResolveReferences(true),
		)

		file, err := astutil.ToFile(node.(*ast.StructLit))
		if err != nil {
			return err
		}

		pkg := &ast.Package{
			Name: ast.NewIdent("info"),
		}
		file.Decls = append([]ast.Decl{pkg}, file.Decls...)

		// fmt.Printf("%#+v\n", file)

		bytes, err := format.Node(
			file,
			format.Simplify(),
		)
		if err != nil {
			return err
		}

		str := string(bytes)
		
		fmt.Println(str)

		*/

	}

	return nil
}
