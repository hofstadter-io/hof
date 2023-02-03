package datamodel

import (
	"fmt"
	"io"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/ast/astutil"
	"cuelang.org/go/cue/format"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

var fieldOptions = []cue.Option{
	cue.Attributes(true),
	cue.Concrete(false),
	cue.Definitions(true),
	cue.Hidden(true),
	cue.Optional(true),
	cue.Docs(true),
}

/*
func (dm *Datamodel) Diff() cue.Value {
	return dm.T.diff
}
*/

func (V *Value) Diff() cue.Value {
	return V.Snapshot.Lense.CurrDiff
}

func (dm *Datamodel) HasDiff() bool {
	return dm.T.hasDiffR()
}

func (V *Value) hasDiffR() bool {
	// load own history
	if V.Hof.Datamodel.History {
		if V.hasDiff() {
			return true
		}
	}

	// recurse if children to load any nested histories
	for _, c := range V.Children {
		if c.T.hasDiffR() {
			return true
		}
	}

	return false
}

func (V *Value) hasDiff() bool {
	return V.Diff().Exists()
}

func (dm *Datamodel) CalcDiffs() error {
	return dm.T.calcDiffR()
}

func (V *Value) calcDiffR() error {
	// load own history
	if V.Hof.Datamodel.History {
		err := V.calcDiff()
		if err != nil {

			return err
		}
	}

	// recurse if children to load any nested histories
	for _, c := range V.Children {
		err := c.T.calcDiffR()
		if err != nil {
			return err
		}
	}

	return nil
}

func (V *Value) calcDiff() error {
	// no history
	if len(V.history) == 0 {
		return nil
	}

	// curr & last vars
	var cv, lv cue.Value
	var cs, ls *Snapshot

	cv = V.Value
	cs = V.Snapshot

	for _, S := range V.history {
		ls, lv = S, S.Value

		// finalize current value
		node := cv.Syntax(
			cue.Final(),
			cue.Docs(true),
			cue.Attributes(true),
			cue.Definitions(true),
			cue.Optional(true),
			cue.Hidden(true),
			cue.Concrete(true),
			cue.ResolveReferences(true),
		)
		cv = cv.Context().BuildExpr(node.(*ast.StructLit))

		// get other value
		label := fmt.Sprintf("ver_%s", ls.Timestamp)
		last := lv.LookupPath(cue.ParsePath(label))

		diff, err := DiffValue(last, cv)
		if err != nil {
			return err
		}
		cs.Lense.CurrDiff = diff

		// we don't want to use 'lv' here
		// because of the ver_TS label
		// which we removed for the diff above anyway
		cv = last
		cs = ls
	}

	return nil
}

func diffDatamodel(dm *Datamodel) error {

	

	/*
	dms, err := LoadDatamodels(args, flgs)
	if err != nil {
		return err
	}

	dms, err = filterDatamodelsByTimestamp(dms, flgs)
	if err != nil {
		return err
	}

	for _, dm := range dms {
		if len(dm.History.Past) == 0 {
			fmt.Printf("%s: no history\n", dm.Name)
		} else {
			past := dm.History.Past[0]
			if flgs.Since != "" {
				past = dm.History.Past[len(dm.History.Past)-1]
			}

			fmt.Printf("// %s -> %s\n%s: ", dm.History.Past[0].Timestamp, dm.Timestamp, dm.Name)
			diff, err := structural.DiffValue(past.Value, dm.Value, nil)
			if err != nil {
				return err
			}
			if !diff.Exists() {
				fmt.Println("{}")
			} else {
				ctx := diff.Context()
				m := ctx.CompileString(orderedMask)
				r, err := structural.MaskValue(m, diff, nil)
				if err != nil {
					return err
				}
				fmt.Println(r)
			}
		}
	}
	*/

	return nil
}

/*
func CalcDatamodelStepwiseDiff(dm *Datamodel) error {
	if dm.History == nil || len(dm.History.Past) == 0 {
		return nil
	}
	past := dm.History.Past

	// loop back through time (checkpoints)
	curr := dm
	for i := 0; i < len(past); i++ {
		// get prev to compare against
		prev := past[i]

		// calculate what needs to be done to prev to get to curr
		diff, err := structural.DiffValue(prev.Value, curr.Value, nil)
		if err != nil {
			return err
		}

		curr.Subsume = prev.Value.Subsume(curr.Value)

		// set changes need to arrive at curr
		curr.Diff = diff
		// update before relooping
		curr = prev
	}
	// TODO(subsume), descend into Models and Fields for diff / subsume for more granular information

	return nil
}
*/

//
//
/////  Everything below is a copy of structural.Diff
//       we have a copy here because we are thinking
//       it might be different or optimized, while
//       not adding complexity to the more general structural.Diff
//       TODO we still need to think about List diff in both places
//

func DiffValue(orig, next cue.Value) (cue.Value, error) {
	r, ok := diffValue(orig, next)
	if !ok {
		return cue.Value{}, nil
	}
	return r, nil
}

func diffValue(orig, next cue.Value) (cue.Value, bool) {

	switch orig.IncompleteKind() {
	case cue.StructKind:
		// fmt.Println("struct", orig, next)
		return diffStruct(orig, next)

	case cue.ListKind:
		// fmt.Println("list", orig, next)
		return diffList(orig, next)

	default:
		// fmt.Println("leaf", orig, next)
		return diffLeaf(orig, next)
	}
}

func diffStruct(orig, next cue.Value) (cue.Value, bool) {
	ctx := orig.Context()
	result := ctx.CompileString("{}")
	rmv := ctx.CompileString("{}")
	add := ctx.CompileString("{}")
	didAdd := false
	didRmv := false

	// first loop over val
	iter, _ := orig.Fields(fieldOptions...)
	for iter.Next() {
		s := iter.Selector()
		p := cue.MakePath(s)
		u := next.LookupPath(p)

		// check that field exists in from. Should we be checking f.Err()?
		if u.Exists() {
			r, ok := diffValue(iter.Value(), u)
			// fmt.Println("r:", r, ok, p)
			if ok {
				result = result.FillPath(p, r)
			}
		} else {
			// remove if orig not in next
			didRmv = true
			rmv = rmv.FillPath(p, iter.Value())
		}
	}

	// add anything in next that is not in orig
	iter, _ = next.Fields(fieldOptions...)
	for iter.Next() {
		s := iter.Selector()
		p := cue.MakePath(s)
		v := orig.LookupPath(p)

		// check that field exists in from. Should we be checking f.Err()?
		if !v.Exists() {
			didAdd = true
			add = add.FillPath(p, iter.Value())
		}
	}

	if didRmv {
		result = result.FillPath(cue.ParsePath("\"-\""), rmv)
	}
	if didAdd {
		result = result.FillPath(cue.ParsePath("\"+\""), add)
	}

	// checks to see if nothing changed
	i := 0
	iter, _ = result.Fields()
	for iter.Next() {
		i++
	}
	if i == 0 {
		return result, false
	}

	return result, true
}

func diffList(orig, next cue.Value) (cue.Value, bool) {
	ctx := orig.Context()
	oi, _ := orig.List()
	ni, _ := next.List()

	result := []cue.Value{}
	for oi.Next() && ni.Next() {
		v, ok := diffValue(oi.Value(), ni.Value())
		if ok {
			result = append(result, v)
		}
	}

	return ctx.NewList(result...), len(result) != 0
}

func diffLeaf(orig, next cue.Value) (cue.Value, bool) {
	// ss := orig.Path().Selectors()
	// lbl := ss[len(ss)-1]

	// check if they are the same
	// by type, concreteness, and unify
	// if same, no need to include in diff
	if orig.IncompleteKind() == next.IncompleteKind() {
		if orig.IsConcrete() == next.IsConcrete() {
			u := orig.Unify(next)
			if u.Err() == nil {
				return cue.Value{}, false
			}
		}
	}

	// need to know if this is a basic lit, so we know if we are changing a concrete value
	ctx := orig.Context()
	ret := ctx.CompileString("{}")
	ret = ret.FillPath(cue.ParsePath("\"-\""), orig)
	ret = ret.FillPath(cue.ParsePath("\"+\""), next)

	// otherwise, we have a diff to create
	/*
	rmv := ctx.CompileString("{}")
	rmv = rmv.FillPath(cue.MakePath(lbl), orig)
	ret = ret.FillPath(cue.ParsePath("\"-\""), rmv)

	add := ctx.CompileString("{}")
	add = add.FillPath(cue.MakePath(lbl), next)
	ret = ret.FillPath(cue.ParsePath("\"+\""), add)
	*/

	return ret, true
}

func (dm *Datamodel) PrintDiff(out io.Writer, dflags flags.DatamodelPflagpole) error {
	return dm.T.printDiffR(out, dflags)
}

func (V *Value) printDiffR(out io.Writer, dflags flags.DatamodelPflagpole) error {
	// load own history
	if V.Hof.Datamodel.History {
		if V.hasDiff() {
			if err := V.printDiff(out, dflags); err != nil {
				return err
			}
		}
	}

	// recurse if children to load any nested histories
	for _, c := range V.Children {
		if err := c.T.printDiffR(out, dflags); err != nil {
			return err
		}
	}

	return nil
}

func (V *Value) printDiff(out io.Writer, dflags flags.DatamodelPflagpole) error {
	name := V.Hof.Label
	p := cue.ParsePath(name)


	d := V.Diff()
	ctx := d.Context()
	val := ctx.CompileString("_")

	val = val.FillPath(p, d)

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
		Name: ast.NewIdent("diff"),
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
	
	fmt.Fprintln(out, str)

	return nil
}
