package cuetils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	// "cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	// "cuelang.org/go/cue/format"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

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
				ToSlash: false,
			})
		}
	}

	s := w.String()
	fmt.Println(s)

}
