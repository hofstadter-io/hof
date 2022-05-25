package ast

import (
	"strings"
)

type Script struct {
	Path    string
	Content string

	Help string
	Args []string

	Lines  []string
	Errors []error

	Phases []*Phase
	Files  map[string]*File
}

func (S *Script) AddError(e error) {
	S.Errors = append(S.Errors, e)
}

func (S *Script) AddPhase(ph *Phase) {
	S.Phases = append(S.Phases, ph)
}

func (S *Script) AddFile(f *File) {
	if S.Files == nil {
		S.Files = make(map[string]*File)
	}
	S.Files[f.Path] = f
}

func (P *Parser) parseScript(S *Script) (*Script, error) {
	// P.logger.Infof("parseScript: %s", S.Path)
	P.script = S

	// split into lines
	S.Lines = strings.Split(S.Content, "\n")

	// parse lines
	err := P.parseLines()
	if err != nil {
		return P.script, err
	}

	return P.script, nil
}

func (P *Parser) parseLines() (err error) {
	// make sure to reset
	P.lineno = -1

	// loop over all lines
Loop:
	P.IncLine()
	for P.lineno < len(P.script.Lines) {
		line := P.script.Lines[P.lineno]
		cleaned := cleanLine(line)

		// check for empty line
		if len(cleaned) == 0 {
			// clear node when we hit an empty line
			P.node = nil
			goto Loop
		}

		// make or add to a node
		if P.node == nil {
			P.node = &NodeBase{
				docLine: P.lineno,
				begLine: P.lineno,
				endLine: P.lineno,
			}
		} else {
			P.node.SetEndLine(P.lineno)
		}

		// do some checks on the first char
		f1 := line[0:1]
		switch f1 {

		// phase, first one should be longest seen
		case "%":
			err = P.parsePhase()
			if err != nil {
				P.script.Errors = append(P.script.Errors, err)
			}
			goto Loop

		// files: before, during
		case "-", "=":
			err = P.parseFile()
			if err != nil {
				P.script.Errors = append(P.script.Errors, err)
			}
			goto Loop

		// doc string
		case "#":
			c := strings.TrimLeft(line, "#")
			P.node.SetBegLine(P.lineno + 1)
			P.node.SetEndLine(P.lineno + 1)
			P.node.AddComment(c)
			goto Loop

		}

		// check for comment on cleaned lined
		switch cleaned[0:1] {
		// disregarded comment
		case "#":
			P.node = nil
			goto Loop
		}

		// check for multiline
		if strings.HasSuffix(cleaned, "\\") {
			goto Loop
		}

		// if we get this far, it's a command or something similar, let's parse it
		err = P.parseCmd()
		if err != nil {
			P.script.Errors = append(P.script.Errors, err)
		}

		goto Loop
	}

	return nil
}
