package ast

import (
	"os"
)

type File struct {
	NodeBase

	// write before teh script starts
	Before  bool

	// normal file stuff
	Path    string
	Mode    os.FileMode
	Chown   string
	Content string
}

func (P *Parser) parseFile() error {
	// grab header line
	line := stripTrailingWhitespace(P.script.Lines[begLine])
	header := strings.Fields(line)

	if len(header) < 3 {
		return NewScriptError("Invlid file header", P.node)
	}

	prefix, fparts, suffix := header[0], header[1:len(header)-1], header[len(header)-1]
	if len(prefix) < 2 || len(suffix) < 2 {
		return NewScriptError("Invlid file header", P.node)
	}

	before := line[0] == '-'
	fchar := line[0]
	if strings.Count(prefix, fchar) != len(prefix) || strings.Count(suffix, fchar) != len(suffix) {
		return NewScriptError("Invlid file header", P.node)
	}

	var fmode os.FileMode
	fmode = 0644
	chown := ""
	fpath := fparts[0]

	// check for filemode/chown
	if len(fparts) > 1 {
		if len(fparts) == 3 {
			if strings.Contains(fparts[2], ":") {
				chown = fparts[2]
			} else {
				fm, err := strconv.ParseUint(fparts[2], 8, 32)
				ts.Check(err)
				fmode = os.FileMode(fm)
			}
		}
		if strings.Contains(fparts[1], ":") {
			chown = fparts[1]
		} else {
			fm, err := strconv.ParseUint(fparts[1], 8, 32)
			ts.Check(err)
			fmode = os.FileMode(fm)
		}
	}

	// create the file
	F := &File{
		NodeBase: P.node,
		Before: before,
		Path: fpath,
		Mode: fmode,
		Chown: chown
	}
	// overwrite self on Parser
	P.node = F

	// Parse File

	return nil
}
