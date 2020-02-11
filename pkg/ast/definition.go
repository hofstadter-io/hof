package ast

import (
	"fmt"
)

type Definitions []*Definition

func (N *Definitions) GetParseInfo() *ParseInfo {
	return nil
}

func (N *Definitions) Parent() ASTNode {
	return nil
}

func (N *Definitions) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

type Definition struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode

	Name   *Token
	Target *TokenPath
	Body []ASTNode

	// Phases filled
	scope Scope
}

func (N *Definition) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *Definition) Parent() ASTNode {
	return N.parent
}

func (N *Definition) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (N *Definition) DefineInScope(name string, node ASTNode) error {
	// This will happen during the
	return fmt.Errorf("Files do not have their own scope! You should not use this function")
}

func (N *Definition) LookupInScope(path []string) (ASTNode, error) {

	return nil, nil
}
