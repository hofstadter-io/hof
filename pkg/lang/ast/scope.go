package ast

import (
	"fmt"
  "sort"
	"strings"
	"unicode"
)

type Scope map[string]ASTNode

type Scoped interface {
	DefineInScope(name string, node ASTNode) error
	LookupInScope(path []string) (ASTNode, error)
}

type SimpleScopedNode struct {
	BaseNode
	scope     Scope
}

func (N *SimpleScopedNode) GetScopeKeys() []string {
	keys := []string{}
	for key, _ := range N.scope {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

func (N *SimpleScopedNode) DefineInScope(name string, node ASTNode) error {
	if N.scope == nil {
		N.scope = Scope{}
	}
	_, ok := N.scope[name]
	if ok {
		return fmt.Errorf("'%s' defined twice", name)
	}
	N.scope[name] = node
	return nil
}

func (N *SimpleScopedNode) LookupInScope(path []string) (ASTNode, error) {
	var err error

	name, rest := path[0], path[1:]
	existing, ok := N.scope[name]
	if ok {
		if len(rest) > 0 {
			return existing.(Scoped).LookupInScope(rest)
		}
		return existing, nil
	}

	// TODO check parent scope

	err = fmt.Errorf("unknown reference to %s", strings.Join(path, "."))
	return nil, err
}

type PrivacyScopedNode struct {
	publicScope Scope
	privateScope Scope
}

func (N *PrivacyScopedNode) GetScopeKeys() []string {
	keys := []string{}
	for key, _ := range N.privateScope {
		keys = append(keys, key)
	}
	for key, _ := range N.publicScope {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

func (N *PrivacyScopedNode) DefineInScope(name string, node ASTNode) error {
	if N.publicScope == nil {
		N.publicScope = Scope{}
		N.privateScope = Scope{}
	}
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

func (N *PrivacyScopedNode) LookupInScope(path []string) (ASTNode, error) {
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

