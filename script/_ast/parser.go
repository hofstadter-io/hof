package ast

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/go-git/go-billy/v5"
	"go.uber.org/zap"
)

type Parser struct {
	config *Config

	logger *zap.SugaredLogger
	fs     billy.Filesystem

	scriptList []*Script
	scriptMap  map[string]*Script


	// parser state
	lineno int
	script *Script
	phase  *Phase
	node   Node
}


func NewParser(config *Config) *Parser {
	P := &Parser{
		config:     config,
		scriptList: []*Script{},
		scriptMap:  map[string]*Script{},
	}

	// Logger
	if config.Logger != nil {
		P.logger = config.Logger
	} else {
		P.initLogger()
	}

	// Filesystem
	if config.Logger != nil {
		P.logger = config.Logger
	}

	return P
}

func (P *Parser) ParseScript(filepath string) (*Script, error) {
	S, err :=  P.setupScript(filepath, nil)
	if err != nil {
		return S, err
	}

	return P.parseScript(S)
}

func (P *Parser) setupScript(filepath string, input interface{}) (S *Script, err error) {
	P.logger.Infow("ParseScript", "filepath", filepath)

	var content string

	if input != nil {
		content, err = P.stringifyContent(input)
	} else {
		content, err = P.readFile(filepath)
	}
	if err != nil {
		return nil, err
	}

	S = &Script{
		Path: filepath,
		Content: content,
	}

	return S, nil
}

// stringifyContent converts a number of options into a string for sure
func (P *Parser) stringifyContent(content interface{}) (string, error) {
	switch c := content.(type) {
	case string:
		return c, nil
	case []byte:
		return string(c), nil
	case *bytes.Buffer:
		if c != nil {
			return c.String(), nil
		}
	case io.Reader:
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, c); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	return "", fmt.Errorf("unsupported content type %T", content)
}

// readFile from P.fs by filepath
func (P *Parser) readFile(filepath string) (string, error) {
	if P.fs == nil {
		b, err := ioutil.ReadFile(filepath)
		return string(b), err
	}

	// some filepath checks
	if filepath == "" {
		return "", fmt.Errorf("readFile: filepath empty")
	}
	if filepath[0:1] != "/" {
		filepath = "/" + filepath
	}

	// Open the file
	F, err := P.fs.Open(filepath)
	if err != nil {
		return "", err
	}

	// Read contents
	b, err := ioutil.ReadAll(F)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func isWhitespace(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

func stripTrailingWhitespace(s string) string {
	i := len(s)
	for i > 0 && isWhitespace(s[i-1]) {
		i--
	}
	return s[0:i]
}

func cleanLine(line string) string {
		line = strings.TrimSpace(line)
		return line
}

func cleanMultiLine(content string) string {
	lines := strings.Split(content, "\n")
	for l, L := range lines {
		L = cleanLine(L)
		L = strings.TrimSuffix(L, "\\")
		lines[l] = L
	}
	return strings.Join(lines, " ")
}
