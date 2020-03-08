package parser

import (
	"strconv"
  "strings"

	"github.com/hofstadter-io/hof/pkg/old-lang/ast"
)

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	switch v.(type) {
	case []interface{}:
		return v.([]interface{})
	default:
		return []interface{}{v}
	}
}

func FileCallback(c *current, pkg, imports, defs  interface{}) (interface{}, error) {
	ret := &ast.File {
    PackageDecl: pkg.(*ast.PackageDecl),
	}
	if imports != nil {
    ret.Imports = imports.([]*ast.Import)
	}
	if defs != nil {
		ret.Definitions = defs.([]ast.ASTNode)
	}

	return ret, nil
}

func PackageDeclCallback(c *current, name interface{}) (interface{}, error) {
  ret := &ast.PackageDecl {
		Name: name.(*ast.Token),
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
  }

  return ret, nil
}

func ImportsCallback(c *current, imports interface{}) (interface{}, error) {
	ret := []*ast.Import {}
	vals := toIfaceSlice(imports)

  for _, i := range vals {
    ret = append(ret, i.(*ast.Import))
  }

	return ret, nil
}

func ImportCallback(c *current, name, path interface{}) (interface{}, error) {
	ret := &ast.Import {
		ImportPath: path.(*ast.String),
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
	}

  if name != nil {
    ret.NameOverride = name.(*ast.Token)
  }

	return ret, nil
}

func DefinitionsCallback(c *current, defs interface{}) (interface{}, error) {
	ret := []ast.ASTNode {}
	vals := toIfaceSlice(defs)

  for _, def := range vals {
    ret = append(ret, def.(ast.ASTNode))
  }

	return ret, nil
}

func TypeDefinitionCallback(c *current, name, paths, body interface{}, isOpen bool) (interface{}, error) {
	ps := toIfaceSlice(paths)
	PS := []*ast.TokenPath{}
	for _, p := range ps {
		if p == nil {
			continue
		}
		switch t := p.(type) {
		/*
		case *ast.TokenPath:
			PS = append(PS, t)
		*/
		case []interface{}:
			p2 := t[0].(*ast.TokenPath)
			PS = append(PS, p2)
		}
	}

	ret := &ast.TypeDefinition {
		Name: name.(*ast.Token),
		Paths: PS,
		Open: isOpen,
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
	}

	if body != nil {
		ret.Body = body.([]ast.ASTNode)
	}

	return ret, nil
}

func DefinitionBodyCallback(c *current, defs interface{}) (interface{}, error) {
	ret := []ast.ASTNode {}

	vals := toIfaceSlice(defs)

  for _, val := range vals {
    ret = append(ret, val.(ast.ASTNode))
  }

	return ret, nil
}

func GeneratorCallback(c *current, id, paths, body interface{}) (interface{}, error) {
	ps := toIfaceSlice(paths)
	PS := []*ast.TokenPath{}
	for _, p := range ps {
		if p == nil {
			continue
		}
		switch t := p.(type) {
		/*
		case *ast.TokenPath:
			PS = append(PS, t)
		*/
		case []interface{}:
			p2 := t[0].(*ast.TokenPath)
			PS = append(PS, p2)
		}
	}

	ret := &ast.GeneratorDef {
		Name: id.(*ast.Token),
		Paths: PS,
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
	}
	if body != nil {
		ret.Body = body.([]ast.ASTNode)
	}

	return ret, nil
}

func EtlDefinitionCallback(c *current, name, args, ret, body interface{}) (interface{}, error) {
	// fmt.Println("EtlDefCallback", name)
	// fmt.Printf("args: %#+v\n", args)
	vals := toIfaceSlice(args)

	es := []*ast.EtlArg{}
  for _, val := range vals {
		switch t := val.(type) {
		case *ast.EtlArg:
			es = append(es, t)
		}
  }
	etl := &ast.EtlDefinition {
		Name: name.(*ast.Token),
		Args: es,
		Return: ret.(*ast.TokenPath),
		Body: body.([]ast.ASTNode),
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
	}

	return etl, nil
}

func EtlArgCallback(c *current, name, path interface{}) (interface{}, error) {
	ret := &ast.EtlArg {
		Name: name.(*ast.Token),
		Path: path.(*ast.TokenPath),
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
	}

	return ret, nil
}

func TypeDefCallback(c *current, id, path, obj interface{}) (interface{}, error) {
	ret := &ast.TypeDef {
		Name: id.(*ast.Token),
		Path: path.(*ast.TokenPath),
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
	}
	if obj != nil {
		objVal := obj.(*ast.Object)
		ret.Extra = objVal
	}

	return ret, nil
}

func FieldTypeCallback(c *current, val interface{}) (interface{}, error) {
	ret := &ast.FieldType {
    Value: val.(ast.ASTNode),
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
	}

	return ret, nil
}

func FieldValueCallback(c *current, val interface{}) (interface{}, error) {
	ret := &ast.FieldValue {
    Value: val.(ast.ASTNode),
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
	}

	return ret, nil
}

func FieldCallback(c *current, id, val interface{}) (interface{}, error) {
	ret := &ast.Field {
    Key: id.(*ast.Token),
    Value: val.(ast.ASTNode),
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
	}

	return ret, nil
}

func ObjectCallback(c *current, fields interface{}) (interface{}, error) {
	vals := toIfaceSlice(fields)

  ret := &ast.Object {
    Fields: make([]*ast.Field, 0, len(vals)),
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
  }

  for _, val := range vals {
      ret.Fields = append(ret.Fields, val.(*ast.Field))
  }

	return ret, nil
}

func ArrayDefCallback(c *current, path interface{}) (interface{}, error) {
	ret := &ast.ArrayDef {
		Path: path.(*ast.TokenPath),
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
	}

	return ret, nil
}

func ArrayCallback(c *current, elems interface{}) (interface{}, error) {
	vals := toIfaceSlice(elems)
	ret := &ast.Array { Elems: make([]ast.ASTNode, 0, len(vals)) }
	for _, val := range vals {
			ret.Elems = append(ret.Elems, val.(ast.ASTNode))
	}
	return ret, nil
}

func PathCallback(c *current, val interface{}) (interface{}, error) {
	text := string(c.text)
  paths := strings.Split(text, ".")
  return &ast.TokenPath {
    Paths: paths,
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
  }, nil
}

func IdCallback(c *current, val interface{}) (interface{}, error) {
  return &ast.Token {
    Value: string(c.text),
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
  }, nil
}

func NameCallback(c *current, val interface{}) (interface{}, error) {
  return &ast.Token {
    Value: string(c.text),
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
  }, nil
}

func TokenCallback(c *current, val interface{}) (interface{}, error) {
  return &ast.Token {
    Value: string(c.text),
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
  }, nil
}

func NumberCallback(c *current) (interface{}, error) {
	// JSON numbers have the same syntax as Go's, and are parseable using
	// strconv.
	val, err := strconv.ParseFloat(string(c.text), 64)
	if err != nil {
		return nil, err
	}

	ret := &ast.Decimal {
    Value: val,
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
  }

	return ret, nil
}

func IntegerCallback(c *current) (interface{}, error) {
	// JSON numbers have the same syntax as Go's, and are parseable using
	val, err := strconv.ParseInt(string(c.text), 10, 64)
	if err != nil {
		return nil, err
	}

	ret := &ast.Integer {
    Value: int(val),
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
  }

	return ret, nil
}

func IntegerDefCallback(c *current) (interface{}, error) {
	ret := &ast.IntegerDef {
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
  }

	return ret, nil
}

func BoolCallback(c *current, val bool) (interface{}, error) {
	ret := &ast.Bool {
    Value: val,
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
  }

	return ret, nil
}

func BoolDefCallback(c *current) (interface{}, error) {
	ret := &ast.BoolDef {
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
  }

	return ret, nil
}

func StringCallback(c *current) (interface{}, error) {
  // TODO : the forward slash (solidus) is not a valid escape in Go, it will
  // fail if there's one in the string
  text, err := strconv.Unquote(string(c.text))
  if err != nil {
      return &ast.Token{}, err
  }

  ret := &ast.String {
    Value: text,
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
  }
  return ret, nil
}

func StringDefCallback(c *current) (interface{}, error) {
	ret := &ast.StringDef {
		BaseNode: ast.BaseNode {
			ParseInfo: ExtractParseInfo(c),
		},
  }

	return ret, nil
}

func ExtractParseInfo(c *current) *ast.ParseInfo {
  return &ast.ParseInfo {
		File: c.globalStore["filename"].(string),
		Line: c.pos.line,
		Col: c.pos.col,
		Offset: c.pos.offset,
		Text: string(c.text),
	}
}
