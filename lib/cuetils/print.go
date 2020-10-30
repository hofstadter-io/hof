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


func CueSyntax(val cue.Value, opts []cue.Option) (ast.Node) {
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

func ValueToSyntaxString(val cue.Value) (string, error) {
	src, err := format.Node(val.Syntax())
	str := string(src)
	return str, err
}

func (CRT *CueRuntime) ParseCueExpr(expr string) (cue.Value, error) {
	inst, err := CRT.CueRuntime.Compile("", expr)
	if err != nil {
		return cue.Value{}, err
	}
	val := inst.Value()
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
		for attrKey, attrVal := range value.Attributes() {
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

func PrintCueError(err error) {

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

	s := w.String()
	fmt.Println(s)

}
func (CR *CueRuntime) PrintCueErrors() {
	for _, err := range CR.CueErrors {
		PrintCueError(err)
	}
}
