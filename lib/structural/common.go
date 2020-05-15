package structural

import (
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/parser"
)

var r cue.Runtime

type pvStruct struct {
	sl *ast.StructLit
}

type pvList struct {
	ll *ast.ListLit
}

type pvExpr struct {
	e *ast.Expr
}

func NewpvStruct() *pvStruct {
	return &pvStruct{
		sl: ast.NewStruct(),
	}
}

func NewpvList() *pvList {
	return &pvList{
		ll: ast.NewList(),
	}
}

func ExprFromValue(v cue.Value) *pvExpr {
	bytes, err := format.Node(v.Syntax())
	if err != nil {
		panic(err)
	}
	expr, err := parser.ParseExpr("", bytes)
	if err != nil {
		panic(err)
	}
	return &pvExpr{
		e: &expr,
	}
}

func (pv pvStruct) Get(key string) *pvStruct {
	for _, d := range pv.sl.Elts {
		df := d.(*ast.Field)
		label, _, _ := ast.LabelName(df.Label)
		if label == key {
			return &pvStruct{
				sl: df.Value.(*ast.StructLit),
			}
		}
	}
	return nil
}

func (pv *pvStruct) Ensure(key string) {
	pvS := pv.Get(key)
	if pvS == nil {
		var e ast.Expr = ast.NewStruct()
		pv.Set(key, pvExpr{
			e: &e,
		})
	}
}

func (pv *pvStruct) Set(key string, expr pvExpr) {
	found := false
	for _, d := range pv.sl.Elts {
		df := d.(*ast.Field)
		label, _, _ := ast.LabelName(df.Label)
		if label == key {
			df.Value = *expr.e
			found = true
			break
		}
	}
	if !found {
		newField := &ast.Field{Label: ast.NewIdent(key), Value: *expr.e}
		pv.sl.Elts = append(pv.sl.Elts, newField)
	}
}

func (pv *pvList) Append(expr pvExpr) {
	pv.ll.Elts = append(pv.ll.Elts, *expr.e)
}

func (pv *pvStruct) ToValue() (*cue.Value, error) {
	i, err := r.CompileExpr(pv.sl)
	if err != nil {
		return nil, err
	}
	v := i.Value()
	if err = v.Err(); err != nil {
		return nil, err
	}

	return &v, nil
}

func (pv *pvList) ToValue() (*cue.Value, error) {
	i, err := r.CompileExpr(pv.ll)
	if err != nil {
		return nil, err
	}
	v := i.Value()
	if err = v.Err(); err != nil {
		return nil, err
	}

	return &v, nil
}
func (pv *pvStruct) ToString() (string, error) {
	v, err := pv.ToValue()
	if err != nil {
		return "", nil
	}
	bytes, err := format.Node(v.Syntax())
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(bytes)), nil
}

func (pv *pvList) ToString() (string, error) {
	v, err := pv.ToValue()
	if err != nil {
		return "", nil
	}
	bytes, err := format.Node(v.Syntax())
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(bytes)), nil
}
func (pv *pvStruct) ToExpr() *pvExpr {
	var e ast.Expr = pv.sl
	return &pvExpr{
		e: &e,
	}
}

func (pv *pvList) ToExpr() *pvExpr {
	var e ast.Expr = pv.ll
	return &pvExpr{
		e: &e,
	}
}

//////////

func isStruct(val cue.Value) bool {
	k := val.Kind()
	return k == cue.StructKind
}

func isList(val cue.Value) bool {
	k := val.Kind()
	return k == cue.ListKind
}

func isBuiltin(val cue.Value) bool {
	k := val.Kind()
	return k == cue.NullKind ||
		k == cue.BoolKind ||
		k == cue.IntKind ||
		k == cue.FloatKind ||
		k == cue.StringKind ||
		k == cue.BytesKind
}
