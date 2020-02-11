package ast

import (
	"fmt"
	"strings"
	"unicode"
)

type TypeDef struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode

	Name *Token
	Path *TokenPath
	Extra *Object

	// Phases filled
	publicScope Scope
	privateScope Scope
}

func (N *TypeDef) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *TypeDef) Parent() ASTNode {
	return N.parent
}

func (N *TypeDef) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (N *TypeDef) DefineInScope(name string, node ASTNode) error {
	// Check first rune to determine public/private
	// Upper Is Public, lower is private
	r := []rune(name)[0]
	if unicode.IsUpper(r) {
		_, ok := N.publicScope[name]
		if ok {
			return fmt.Errorf("'%s' defined twice", name)
			// return fmt.Errorf("'%s' defined twice\n - %s\n - %s\n", name, existing.GetParseInfo(), node.GetParseInfo())
		}
		N.publicScope[name] = node
	} else {
		_, ok := N.privateScope[name]
		if ok {
			return fmt.Errorf("'%s' defined twice", name)
			// return fmt.Errorf("'%s' defined twice\n - %s\n - %s\n", name, existing.GetParseInfo(), node.GetParseInfo())
		}
		N.privateScope[name] = node
	}
	return nil
}

func (N *TypeDef) LookupInScope(path []string) (ASTNode, error) {
	var err error

	name, rest := path[0], path[1:]
	// Check first rune to determine public/private
	// Upper Is Public, lower is private
	r := []rune(name)[0]
	if unicode.IsUpper(r) {
		existing, ok := N.publicScope[name]
		if ok {
			if len(rest) > 0 {
				return existing.(Scoped).LookupInScope(rest)
			}
			return existing, nil
		}
	} else {
		existing, ok := N.privateScope[name]
		if ok {
			if len(rest) > 0 {
				return existing.(Scoped).LookupInScope(rest)
			}
			return existing, nil
		}
	}

	err = fmt.Errorf("unknown reference to %s", strings.Join(path, "."))
	return nil, err
}

