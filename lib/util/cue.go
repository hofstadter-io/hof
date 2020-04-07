package util

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"cuelang.org/go/cue/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

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
		errors.Print(w, err, &errors.Config{
			Format:  format,
			Cwd:     cwd,
			ToSlash: false,
		})
		s := w.String()
		fmt.Println(s)

}

