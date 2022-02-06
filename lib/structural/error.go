package structural

import (
	"bytes"
	"io"
	"os"
	"strings"

	"cuelang.org/go/cue"
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

func FormatCueError(err error) string {

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
	return s
}

func getErrorAttrMsg(val cue.Value) (string, bool) {
	msg, has := "", false
	attr := val.Attribute("error")
	if attr.Err() == nil {
		has = true
		if attr.NumArgs() > 0 {
			m, _ := attr.String(0)
			if m != "" {
				msg = m
			}
		}
	}
	return msg, has
}
