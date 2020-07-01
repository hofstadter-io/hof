package ast

import (
	"strings"
)

func (P *Parser) readSingleQuoteString(input string) (token, rest string, err error) {
	// P.logger.Errorf("'xQ: %q", input)
	if len(input) == 0 {
		return "", "", NewScriptError("readSingleQuoteString arg has zero length?", P.node, nil)
	}
	if input[0] != '\'' {
		return "", "", NewScriptError("readSingleQuoteString arg missing leading quote", P.node, nil)
	}

	if len(input) >= 3 && input[0:3] == "'''" {
		return P.readTripleQuoteString("'''", input)
	}

	p := 1
	for p < len(input) {
		if input[p] == '\'' && input[p-1] != '\\' {
			return input[0:p+1], input[p+1:], nil
		}
		p++
	}

	return "", "", NewScriptError("readSingleQuoteString string missing final quote", P.node, nil)
}

func (P *Parser) readDoubleQuoteString(input string) (token, rest string, err error) {
	// P.logger.Errorf("\"xQ: %q", input)
	if len(input) == 0 {
		return "", "", NewScriptError("readDoubleQuoteString arg has zero length?", P.node, nil)
	}

	if input[0] != '"' {
		return "", "", NewScriptError("readDoubleQuoteString arg missing leading quote", P.node, nil)
	}

	if len(input) >= 3 && input[0:3] == `"""` {
		return P.readTripleQuoteString(`"""`, input)
	}

	p := 1
	for p < len(input) {
		if input[p] == '"' && input[p-1] != '\\' {
			return input[0:p+1], input[p+1:], nil
		}
		p++
	}

	return "", "", NewScriptError("readDoubleQuoteString string missing final quote", P.node, nil)
}

func (P *Parser) readBacktickString(input string) (token, rest string, err error) {
	// P.logger.Errorf("`xQ: %q", input)
	if len(input) == 0 {
		return "", "", NewScriptError("readBacktickString arg has zero length?", P.node, nil)
	}

	if input[0] != '`' {
		return "", "", NewScriptError("readBacktickString arg missing leading quote", P.node, nil)
	}

	if len(input) >= 3 && input[0:3] == "```" {
		return P.readTripleQuoteString("```", input)
	}

	p := 1
	for p < len(input) {
		if input[p] == '`' && input[p-1] != '\\' {
			return input[0:p+1], input[p+1:], nil
		}
		p++
	}

	return "", "", NewScriptError("readBacktickString string missing final quote", P.node, nil)
}

func (P *Parser) readTripleQuoteString(quote, input string) (token, rest string, err error) {
	// P.logger.Errorf("3xQ: %q %q", quote, input)
	if input[0:3] != quote {
		return "", "", NewScriptError("readTripleQuoteString arg missing leading quote", P.node, nil)
	}

	// handle full text a special way, input should really be a single line unless nested
	if strings.Contains(input, "\n") {
		lines := strings.Split(input, "\n")
		content := lines[0]
		for _, line := range lines {
			if strings.HasPrefix(line, quote) {
				content += quote + "\n"
				return content, "", nil
			}
			content += line + "\n"
		}
		return content, "", NewScriptError("readTripleQuoteString string missing final quote", P.node, nil)
	}

	// now reading triple quote content
	content := input + "\n"
	P.IncLine()
	for !P.EOF() {
		line := P.script.Lines[P.lineno]

		// strop when we find the end of the quote
		if strings.HasPrefix(line, quote) {
			content += quote + "\n"
			// see about continuing this command
			rest = strings.TrimPrefix(line, quote)
			rest := cleanLine(rest)
			return content, rest, nil
		}
		content += line + "\n"
		P.IncLine()
	}

	return content, "", NewScriptError("readTripleQuoteString string missing final quote", P.node, nil)
}
