package ast

import (
	"fmt"
	"strings"
	"unicode"
)

type CmdExpect int

const (
	None CmdExpect = iota
	Pass
	Fail
	Skip
)

type Cmd struct {
	NodeBase

	Exp  CmdExpect
	Cmd  string
	Args []string
	Cmds []*Cmd
	Bg     bool
	BgName string
}

func (P *Parser) parseCmd() error {
	S := P.script
	N := P.node

	// fixup the lines, removing space and removing multiline suffix
	lines := S.Lines[N.BegLine():N.EndLine()+1]
	for i, l := range lines {
		l = strings.TrimSpace(l)
		l = strings.TrimSuffix(l, "\\")
		lines[i] = l
	}

	// Build a single line from the slice
	line := strings.Join(lines, "")

	// parse the line(s) into a command
	cmd, err := P.parseCmdLine(line)
	if err != nil {
		return err
	}

	// add the command to things
	P.AppendNode(cmd)
	return nil
}

// parses a line into a command
func (P *Parser) parseCmdLine(line string) (*Cmd, error) {
	P.logger.Debugf("parseCmdLine: %q", line)

	cmd := &Cmd{
		NodeBase: P.node.CloneNodeBase(),
		Exp: Pass,
	}

	for len(line) > 0 {
		token, rest, err := P.nextToken(line)
		// P.logger.Warnf("Token: %q", token)
		if err != nil {
			return cmd, NewScriptError("While parsing cmd", cmd, err)
		}
		if len(token) == 0 {
			return cmd, nil
		}
		line = rest

		switch token[0] {

			case '!':
				cmd.Exp = Fail
			case '?':
				cmd.Exp = Skip

			case '&':
				flds := strings.SplitN(token, ":", 2)
				cmd.Bg = true
				cmd.BgName = flds[1]

			// single (3x?) quoted arg
			case '\'':
				// P.logger.Infof("   single: %q", token)
				cmd.Args = append(cmd.Args, token)

			// double (3x?) quoted arg
			case '"':
				// P.logger.Infof("   double: %q", token)
				cmd.Args = append(cmd.Args, token)

			// backtick (3x?) quoted exec'n arg
			case '`':
				// P.logger.Infof("   backtk: %q", token)
				// find the actual subcommand
				i := 1
				subtoken := token[i:len(token)-i]
				if strings.HasPrefix(token, "```") {
					i = 3
					subtoken = token[i:len(token)-(i+1)]
				}

				// parse the subcommand
				subcmd, err := P.parseCmdLine(subtoken)
				if err != nil {
					return subcmd, err
				}

				// add to top level command and continue
				cmd.Args = append(cmd.Args, fmt.Sprintf("$HOF_INLINE_CMD_%d", len(cmd.Cmds)))
				cmd.Cmds = append(cmd.Cmds, subcmd)

			default:
				// first time through should set the command
				if cmd.Cmd == "" {
					if token == "wait" {
						subcmd, err := P.parseCmdLine(rest)
						if err != nil {
							return subcmd, err
						}
						cmd.Cmds = append(cmd.Cmds, subcmd)
					}
					cmd.Cmd = token
					break
				} else {
					// all other times are ags
					cmd.Args = append(cmd.Args, token)
				}
		}

	}

	return cmd, nil
}

func (P *Parser) nextToken(input string) (token, rest string, err error) {
	var (
		p int
		r rune
	)

	// consume spaces
	for i, c := range input {
		// break at first non-space
		if !unicode.IsSpace(c) {
			p = i
			r = c
			break
		}
	}

	switch r {
	case '!', '?', '=':
		return string(r), input[p+1:], nil

	case '&':
		return P.readBackgroundToken(input[p:])

	case '\'':
		return P.readSingleQuoteString(input[p:])
	case '"':
		return P.readDoubleQuoteString(input[p:])
	case '`':
		return P.readBacktickString(input[p:])

	// TODO, multiline strings
	case '#':
		return readHashtag(input[p:])
	}

	// read until we hit a space or EOL
	for i, c := range input[p:] {
		if c == '=' {
			p = i - 1
			break
		}
		if unicode.IsSpace(c) {
			p = i
			break
		}
		token += string(c)
		p = i
	}

	return token, input[p+1:], nil
}

func (P *Parser) readBackgroundToken(input string) (token, rest string, err error) {
	if input[0] != '&' {
		return "", "", fmt.Errorf("readBackgroundToken arg missing leading &")
	}

	for _, c := range input[1:] {
		if unicode.IsSpace(c) {
			return string(input[0]), input[2:], nil
		}
		if c == ':' {
			name, rest, err := P.nextToken(input[2:])
			return input[0:2]+name, rest, err
		}

		return "", "", fmt.Errorf("Background command missing space or :name after &")
	}

		return "", "", fmt.Errorf("readBackgroundToken shouldn't get here")
}

func readHashtag(input string) (token, rest string, err error) {
	if input[0] != '#' {
		return "", "", fmt.Errorf("readBackgroundToken arg missing leading &")
	}

	p := 1
	for _, c := range input[1:] {
		p++
		if c == '\n' {
			break
		}
	}

	return "", input[p:], nil
}
