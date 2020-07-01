package ast

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

type File struct {
	NodeBase

	// write out before the script starts
	Before  bool
	// trim the end
	TrimEnd bool

	// normal file stuff
	Path    string
	Mode    os.FileMode
	Chown   string

	// content markers
	BegContent int
	EndContent int
}

func (P *Parser) parseFile() error {
	N := P.node
	S := P.script

	// grab header line
	line := stripTrailingWhitespace(S.Lines[N.BegLine()])
	header := strings.Fields(line)
	P.logger.Debugf("parseFile: %q %v", line, header)

	if len(header) < 3 {
		return NewScriptError("Invalid file header format", N, nil)
	}

	prefix, fparts, suffix := header[0], header[1:len(header)-1], header[len(header)-1]
	if len(prefix) < 2 || len(suffix) < 2 {
		return NewScriptError("Invalid file header delims", N, nil)
	}

	before := line[0] == '-'
	fchar := line[0:1]
	if strings.Count(prefix, fchar) != len(prefix) || strings.Count(suffix, fchar) != len(suffix) {
		return NewScriptError("Invalid file header delims", N, nil)
	}

	fpath := fparts[0]
	fargs := fparts[1:]

	// create the file
	F := &File{
		NodeBase: P.node.CloneNodeBase(),
		Before: before,
		Path: fpath,
		Mode: 0644,
		Chown: "",
		// set to header line so we can compare later
		BegContent: P.lineno,
		EndContent: P.lineno,
	}

	// loop over args
	for _, arg := range fargs {
		// check for trimming
		if arg == "trim" || arg == "//" {
			F.TrimEnd = true
			continue
		}

		// check for filemode syntax
		if b, err := regexp.MatchString("[[:digit:]]+", arg); b || err != nil {
			// this should never happen at runtime
			if err != nil {
				panic(err)
			}

			fm, cerr := strconv.ParseUint(arg, 8, 32)
			if cerr != nil {
				return NewScriptError("Invalid file header filemode: " + arg, F, err)
			}
			F.Mode = os.FileMode(fm)
			continue
		}
		if b, err := regexp.MatchString("[-xwr]+", arg  ); b || err != nil {
			// this should never happen at runtime
			if err != nil {
				panic(err)
			}

			return NewScriptError("char based filemode is not supported yet", F, err)
		}

		// otherwise assume chown
		F.Chown = arg
	}

	// overwrite self on Parser
	P.node = Node(F)

	//
	//////   Parse File
	//

	// inc lineno to first line of content
	P.lineno++
	// start consuming lines
	for !P.EOF() {
		l := S.Lines[P.lineno]

		// does the line start with the file header prefix? (it should match to indicate transition)
		if strings.HasPrefix(l, prefix){
			flds := strings.Fields(l)

			// is the first field the same? and there are at least 2 fields? let's make a transition
			if len(flds[0]) >= len(prefix) && len(flds) > 2 {
				// if we have "<prefix> end .*", consume the line
				// otherwise we are assuming the start of the next file and want to keep it
				if strings.TrimSpace(strings.ToLower(flds[1])) == "end" {
					P.IncLine()
				}

				// finalize the file and return to main loop
				S.AddFile(F)
				P.AppendNode(F)
				return nil
			}
		}

		// ok, no transition, so accum, inc, and try the next line
		// accume and set endings
		F.SetEndLine(P.lineno)
		F.EndContent = P.lineno

		// set initial the first time around
		// we may never get here if a file has no content
		if F.BegContent == F.BegLine() {
			F.BegContent = P.lineno
		}

		// inc and loop around
		P.IncLine()
	}

	// finalize this file
	S.AddFile(F)
	P.AppendNode(F)

	return nil
}
