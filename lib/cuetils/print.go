package cuetils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/format"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func PrintCue(val cue.Value) (string, error) {
	syn := val.Syntax(
		cue.Final(),
		cue.ResolveReferences(true),
		cue.Definitions(true),
		cue.Hidden(true),
		cue.Optional(true),
		cue.Attributes(true),
		cue.Docs(true),
	)

	bs, err := format.Node(syn)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func FormatCue(val cue.Value) (string, error) {
	syn := val.Syntax(
		cue.Final(),
		cue.ResolveReferences(true),
		cue.Concrete(true),
		cue.Definitions(true),
		cue.Hidden(true),
		cue.Optional(true),
		cue.Attributes(true),
		cue.Docs(true),
	)

	bs, err := format.Node(syn)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func CueSyntax(val cue.Value, opts []cue.Option) ast.Node {
	if len(opts) > 0 {
		return val.Syntax(opts...)
	}
	return val.Syntax(
		cue.Attributes(true),
		cue.Concrete(false),
		cue.Definitions(true),
		cue.Docs(true),
		cue.Hidden(true),
		cue.Optional(true),
		cue.ResolveReferences(true),
	)
}

func PrintCueValue(val cue.Value) (string, error) {
	node := val.Syntax(
		cue.Attributes(true),
		cue.Concrete(false),
		cue.Definitions(true),
		cue.Docs(true),
		cue.Hidden(true),
		cue.Final(),
		cue.Optional(false),
		cue.ResolveReferences(true),
	)

	bytes, err := format.Node(
		node,
		format.TabIndent(false),
		format.UseSpaces(2),
		format.Simplify(),
	)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ValueToSyntaxString(val cue.Value, opts ...cue.Option) (string, error) {
	src, err := format.Node(val.Syntax(opts...))
	str := string(src)
	return str, err
}

func (CRT *CueRuntime) ParseCueExpr(expr string) (cue.Value, error) {
	val := CRT.CueContext.CompileString(expr)
	if val.Err() != nil {
		return val, val.Err()
	}
	return val, nil
}

func (CRT *CueRuntime) PrintValue() error {
	// Get top level struct from cuelang
	S, err := CRT.CueValue.Struct()
	if err != nil {
		return err
	}

	iter := S.Fields()
	for iter.Next() {

		label := iter.Label()
		value := iter.Value()
		fmt.Println("  -", label, value)
		for attrKey, attrVal := range value.Attributes(cue.ValueAttr) {
			fmt.Println("  --", attrKey)
			for i := 0; i < 5; i++ {
				str, err := attrVal.String(i)
				if err != nil {
					break
				}
				fmt.Println("  ---", str)
			}
		}
	}

	return nil
}

func getLang() language.Tag {
	loc := os.Getenv("LC_ALL")
	if loc == "" {
		loc = os.Getenv("LANG")
	}
	loc = strings.Split(loc, ".")[0]
	return language.Make(loc)
}

func CueErrorToString(err error) string {

	p := message.NewPrinter(getLang())
	format := func(w io.Writer, format string, args ...interface{}) {
		p.Fprintf(w, format, args...)
	}
	cwd, _ := os.Getwd()
	w := &bytes.Buffer{}

	for _, e := range errors.Errors(err) {
		for _, e2 := range errors.Errors(e) {
			errors.Print(w, e2, &errors.Config{
				Format:  format,
				Cwd:     cwd,
				ToSlash: true,
			})
		}
	}

	return w.String()
}

func ExpandCueError(err error) error {
	s := CueErrorToString(err)
	return fmt.Errorf(s)
}

func PrintCueError(err error) {
	s := CueErrorToString(err)
	fmt.Println(s)
}

func (CR *CueRuntime) PrintCueErrors() {
	for _, err := range CR.CueErrors {
		PrintCueError(err)
	}
}
