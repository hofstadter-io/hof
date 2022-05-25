package ast

import (
	"github.com/go-git/go-billy/v5"
	"github.com/kr/pretty"
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

func (P *Parser) IncLine() {
	// update some other internal trackers
	if P.phase != nil {
		P.phase.SetEndLine(P.lineno)
	}
	if P.node != nil {
		P.node.SetEndLine(P.lineno)
	}

	// update Parsers lineno
	P.lineno++
}

func (P *Parser) AppendNode(n Node) {
	// this only happens when the user hasn't started with a phase
	if P.phase == nil {
		P.phase = &Phase{
			NodeBase: NodeBase{
				name:    "unnamed phase",
				docLine: 0,
				begLine: 1,
				endLine: P.lineno,
			},
			level: 0,
			title: "unnamed phase",
		}
		P.script.AddPhase(P.phase)
	}

	// append and clear
	P.phase.AppendNode(n)
	P.node = nil
}

func (P *Parser) ParseScript(filepath string) (*Script, error) {
	S, err := P.setupScript(filepath, nil)
	if err != nil {
		return S, err
	}

	S, err = P.parseScript(S)

	P.logger.Infof("Script AST:\n%# v\n", pretty.Formatter(S))

	return S, err
}

func (P *Parser) setupScript(filepath string, input interface{}) (S *Script, err error) {
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
		Path:    filepath,
		Content: content,
	}

	return S, nil
}
