package ast

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
	"github.com/kr/pretty"
)

func ValueFromData(input interface{}) (ASTNode, error) {
	switch data := input.(type) {
	case map[string]interface{}:
		_, nameOk := data["name"]
		_, typeOk := data["type"]

		if nameOk && typeOk {

			var N TypeDecl
			ret, err := N.FromData(input)
			if err != nil {
				return nil, err
			}
			return ret, nil

		} else {

			var N Object
			ret, err := N.FromData(input)
			if err != nil {
				return nil, err
			}
			return ret, nil

		}


	case []interface{}:
		var N Array
		ret, err := N.FromData(input)
		if err != nil {
			return nil, err
		}
		return ret, nil

	case string:
		var N Token
		ret, err := N.FromData(input)
		if err != nil {
			return nil, err
		}
		return ret, nil

	case int:
		var N Integer
		ret, err := N.FromData(input)
		if err != nil {
			return nil, err
		}
		return ret, nil

	case float64:
		var N Decimal
		ret, err := N.FromData(input)
		if err != nil {
			return nil, err
		}
		return ret, nil

	case bool:
		var N Bool 
		ret, err := N.FromData(input)
		if err != nil {
			return nil, err
		}
		return ret, nil

	default:
		return nil, errors.New("Unknown type in ValueFromData: " + fmt.Sprintf("%# v", pretty.Formatter(input)) )
	}
}

func (N HofFile) FromData(input interface{}) (ASTNode, error) {
	switch data := input.(type) {
	case []interface{}:
		for _, elem := range data {
			var def Definition
			tmp, err := def.FromData(elem)
			if err != nil {
				return N, err
			}
			N.Definitions = append(N.Definitions, tmp.(Definition))
		}

	case map[string]interface{}:
		var def Definition
		tmp, err := def.FromData(data)
		if err != nil {
			return N, err
		}
		N.Definitions = append(N.Definitions, tmp.(Definition))
	}

	return N, nil
}

func (N Definition) FromData(input interface{}) (ASTNode, error) {
	data := input.(map[string]interface{})

	// should only be one key here
	var dsl_type string
	var dsl_content interface{}
	for key, val := range data {
		dsl_type = key
		dsl_content = val
		break
	}

	dsl_map := dsl_content.(map[string]interface{})

	N.DSL = Token { Value: dsl_type }
	N.Name = Token { Value: dsl_map["name"].(string) }
	delete(dsl_map, "name")

	for key, val := range dsl_map {
		if key == "fields" {
			fields := val.([]interface{})
			for _, field := range fields {
				var decl TypeDecl
				tmp, err := decl.FromData(field)
				if err != nil {
					return N, err
				}
				N.Body = append(N.Body, tmp.(TypeDecl))
			}

		} else {

			// figure out what to parse / build next
			value, err := ValueFromData(val)
			if err != nil {
				return N, err
			}
			field := Field {
				Key: Token{ Value: key },
				Value: value,
			}
			N.Body = append(N.Body, field)
		}
	}

	return N, nil
}

func (N TypeDecl) FromData(input interface{}) (ASTNode, error) {
	data := input.(map[string]interface{})

	var name, typ Token

	tmpName, err := name.FromData(data["name"])
	if err != nil {
		return N, err
	}

	tmpType, err := typ.FromData(data["type"])
	if err != nil {
		return N, err
	}

	N.Name = tmpName.(Token)
	N.Type = tmpType.(Token)

	delete(data, "name")
	delete(data, "type")

	if len(data) == 0 {
		return N, nil
	}

	var obj Object
	tmp, err := obj.FromData(data)
	if err != nil {
		return N, err
	}

	extra := tmp.(Object)
	N.Extra = &extra


	return N, nil
}

func (N Object) FromData(input interface{}) (ASTNode, error) {
	data := input.(map[string]interface{})

	for key,val := range data {
		// determine val's type here
		iVal, err := ValueFromData(val)
		if err != nil {
			return N, err
		}
		field := Field {
			Key: Token { Value: key },
			Value: iVal,
		}
		N.Fields = append(N.Fields, field)
	}

	return N, nil
}

func (N Field) FromData(input interface{}) (ASTNode, error) {

	return N, nil
}

func (N Array) FromData(input interface{}) (ASTNode, error) {
	data := input.([]interface{})

	for _, val := range data {
		// determine val's type here
		iVal, err := ValueFromData(val)
		if err != nil {
			return N, err
		}
		N.Elems = append(N.Elems, iVal)
	}

	return N, nil
}

func (N PathExpr) FromData(input interface{}) (ASTNode, error) {

	return N, nil
}

func (N TokenPath) FromData(input interface{}) (ASTNode, error) {

	return N, nil
}

func (N BracePath) FromData(input interface{}) (ASTNode, error) {

	return N, nil
}

func (N RangeExpr) FromData(input interface{}) (ASTNode, error) {

	return N, nil
}

func (N Token) FromData(input interface{}) (ASTNode, error) {
	N.Value = input.(string)
	return N, nil
}

func (N Integer) FromData(input interface{}) (ASTNode, error) {
	N.Value = input.(int)
	return N, nil
}

func (N Decimal) FromData(input interface{}) (ASTNode, error) {
	v := input.(float64)
	if v == math.Trunc(v) {
		ret := Integer {
			Value: int(v),
		}
		return ret, nil
	}

	N.Value = v
	return N, nil
}

func (N Bool) FromData(input interface{}) (ASTNode, error) {
	N.Value = input.(bool)
	return N, nil
}

