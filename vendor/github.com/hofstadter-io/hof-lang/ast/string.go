package ast

import (
	"fmt"
	// "github.com/kr/pretty"
)

func (N HofFile) String(indent string) (string, error) {
	ret := ""
	var part string
	var err error

	for _, c := range N.Definitions {
		part, err = c.String(indent)
		if err != nil {
			return "", err
		}
		ret += part
	}

	return ret, nil
}

func (N Definition) String(indent string) (string, error) {
	ret := ""
	var part string
	var err error

	ret += fmt.Sprintf("%sdef %s %s {\n\n", indent, N.Name.Value, N.DSL.Value)
	for _, c := range N.Body {
		part, err = c.String(indent + "  ")
		if err != nil {
			return "", err
		}
		if part[len(part) - 1:] != "\n" {
			part += "\n"
		}
		ret += part + "\n"
	}
	ret += fmt.Sprintf("%s}\n\n", indent)

	return ret, nil
}

func (N TypeDecl) String(indent string) (string, error) {
	ret := ""
	var part string
	var err error

	ret +=  fmt.Sprintf("%s%s %s ", indent, N.Name.Value, N.Type.Value)
	if N.Extra != nil {
		part, err = N.Extra.String(indent)
		if err != nil {
			return "", err
		}
		ret += part
	}
	return ret, nil
}

func (N Object) String(indent string) (string, error) {
	ret := ""
	var part string
	var err error

	ret += fmt.Sprintf("{")
	if len(N.Fields) == 0 {
		ret += "}\n"
		return ret, nil
	} else {
		ret += "\n"
	}
	for _, c := range N.Fields {
		part, err = c.String(indent + "  ")
		if err != nil {
			return "", err
		}
		ret += part + "\n"
	}
	ret += fmt.Sprintf("%s}", indent)

	return ret, nil
}

func (N Field) String(indent string) (string, error) {
	ret := ""
	var part string
	var err error

	ret += fmt.Sprintf("%s%s: ", indent, N.Key.Value)

	part, err = N.Value.String(indent)
	if err != nil {
		return "", err
	}
	ret += part

	ret += fmt.Sprintf("%s", indent)

	return ret, nil
}

func (N Array) String(indent string) (string, error) {
	ret := ""
	var part string
	var err error

	ret += fmt.Sprintf("[")
	if len(N.Elems) == 0 {
		ret += "]"
		return ret, nil
	} else {
		ret += "\n"
	}
	for _, c := range N.Elems {
		part, err = c.String(indent + "  ")
		if err != nil {
			return "", err
		}

		switch c.(type) {
		case Object, Token, Integer, Decimal, Bool:
			ret += indent + "  "
		}

		ret += part + ",\n"
	}
	ret += fmt.Sprintf("%s]", indent)

	return ret, nil
}

func (N PathExpr) String(indent string) (string, error) {
	ret := ""
	var part string
	var err error

	for _, c := range N.PathList {
		part, err = c.String(indent + "  ")
		if err != nil {
			return "", err
		}
		ret += part
	}

	return ret, nil
}

func (N TokenPath) String(indent string) (string, error) {
	ret := ""

	ret += fmt.Sprintf(".%s", N.Value)

	return ret, nil
}

func (N BracePath) String(indent string) (string, error) {
	ret := ""
	var part string
	var err error

	ret += fmt.Sprintf(".{")
	for _, c := range N.Exprs {
		part, err = c.String(indent + "  ")
		if err != nil {
			return "", err
		}
		ret += part
	}
	fmt.Sprintf("}")

	return ret, nil
}

func (N RangeExpr) String(indent string) (string, error) {
	var ret string

	if !N.Range {
		ret = fmt.Sprintf(".[%d]", indent, N.Low)
		return ret, nil
	}

	low, high := "", ""
	if N.Low >= 0 {
		low = fmt.Sprintf("%d", N.Low)
	}
	if N.High >= 0 {
		high = fmt.Sprintf("%d", N.High)
	}

	ret = fmt.Sprintf(".[%s:%s]", indent, low, high)

	return ret, nil
}

func (N Token) String(indent string) (string, error) {
	return fmt.Sprintf("%q", N.Value), nil
}

func (N Integer) String(indent string) (string, error) {
	return fmt.Sprintf("%d", N.Value), nil
}

func (N Decimal) String(indent string) (string, error) {
	return fmt.Sprintf("%f", N.Value), nil
}

func (N Bool) String(indent string) (string, error) {
	return fmt.Sprintf("%v", N.Value), nil
}
