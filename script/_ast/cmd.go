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
	Bg     bool
	BgName string
}

func (P *Parser) parseCmd() error {
	S := P.script
	N := P.node

	lines := S.Lines[N.BegLine():N.EndLine()+1]
	for i, l := range lines {
		l = strings.TrimSpace(l)
		l = strings.TrimSuffix(l, "\\")
		lines[i] = l
	}

	line := strings.Join(lines, "")
	P.logger.Debugf("parseCmd: %q", line)

	cmd := &Cmd{
		NodeBase: P.node.CloneNodeBase(),
		Exp: Pass,
	}

	for len(line) > 0 {
		token, rest, err := nextToken(line)
		if err != nil {
			return NewScriptError("While parsing cmd", N, err)
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

			default:
				if cmd.Cmd == "" {
					cmd.Cmd = token
				} else {
					cmd.Args = append(cmd.Args, token)
				}
		}

	}

	P.AppendNode(cmd)

	return nil
}

func nextToken(input string) (token, rest string, err error) {
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
	case '!', '?':
		return string(r), input[p+1:], nil

	case '&':
		return readBackgroundToken(input[p:])

	case '\'':
		return readSingleQuoteString(input[p:])
	case '"':
		return readDoubleQuoteString(input[p:])

	// TODO, multiline strings
	}

	// read until we hit a space or EOL
	for i, c := range input[p:] {
		if unicode.IsSpace(c) {
			p = i
			break
		}
		token += string(c)
	}

	return token, input[p:], nil
}

func readSingleQuoteString(input string) (token, rest string, err error) {
	if input[0] != '\'' {
		return "", "", fmt.Errorf("readSingleQuoteString arg missing leading quote")
	}

	p := 1
	for p < len(input) {
		if input[p] == '\'' && input[p-1] != '\\' {
			return input[0:p+1], input[p+1:], nil
		}
		p++
	}

	return "", "", fmt.Errorf("readSingleQuoteString string missing final quote")
}

func readDoubleQuoteString(input string) (token, rest string, err error) {
	if input[0] != '"' {
		return "", "", fmt.Errorf("readDoubleQuoteString arg missing leading quote")
	}

	p := 1
	for p < len(input) {
		if input[p] == '"' && input[p-1] != '\\' {
			return input[0:p+1], input[p+1:], nil
		}
		p++
	}

	return "", "", fmt.Errorf("readDoubleQuoteString string missing final quote")
}

func readBackgroundToken(input string) (token, rest string, err error) {
	if input[0] != '&' {
		return "", "", fmt.Errorf("readBackgroundToken arg missing leading &")
	}

	for _, c := range input[1:] {
		if unicode.IsSpace(c) {
			return string(input[0]), input[2:], nil
		}
		if c == ':' {
			name, rest, err := nextToken(input[2:])
			return input[0:1]+name, rest, err
		}

		return "", "", fmt.Errorf("Background command missing space or :name after &")
	}

		return "", "", fmt.Errorf("readBackgroundToken shouldn't get here")
}
