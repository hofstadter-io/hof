package yagu

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/format"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	// Global cue runtime
	CueRuntime cue.Runtime
)

func PrintCueInstance(i *cue.Instance) error {
	bytes, err := format.Node(i.Value().Syntax())
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
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
