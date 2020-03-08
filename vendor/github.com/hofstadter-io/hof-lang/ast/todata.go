package ast

import (
	"fmt"
	"github.com/pkg/errors"
)

func (N HofFile) ToData() (interface{}, error) {
	data := make([]interface{}, 0, len(N.Definitions))

	for _, def := range N.Definitions {
		defData, err := def.ToData()
		if err != nil {
			return nil, err
		}
		data = append(data, defData)
	}

	return data, nil
}

func (N Definition) ToData() (interface{}, error) {
	data := map[string]interface{}{}

	for _, elem := range N.Body {
		elemData, err := elem.ToData()
		if err != nil {
			return nil, err
		}

		switch elem.(type) {
		case TypeDecl:
			// special case  for types
			fs, ok := data["fields"]
			if ok {
				fss := fs.([]interface{})
				data["fields"] = append(fss, elemData)
			} else {
				data["fields"] =  []interface{}{ elemData }
			}

		default:
			elemMap := elemData.(map[string]interface{})
			for key, val := range elemMap  {
				// should merge here? or error on confict?
				currVal, ok := data[key]
				if ok {
					if key == "fields" {
						fss := currVal.([]interface{})
						elemSlice := elemMap["fields"].([]interface{})
						data["fields"] = append(fss, elemSlice...)
						continue
					}

					fmt.Println(data[key])
					return nil, errors.New("'"+key+"' already defined")
				}
				data[key] = val
			}
		}

	}


	data["name"] = N.Name.Value
	dsl := map[string]interface{}{
		N.DSL.Value: data,
	}

	return dsl, nil
}

func (N TypeDecl) ToData() (interface{}, error) {
	data := map[string]interface{}{
		"type": N.Type.Value,
		"name": N.Name.Value,
	}

	if N.Extra != nil {
		extraData, err := N.Extra.ToData()
		if err != nil {
			return nil, err
		}
		extraMap := extraData.(map[string]interface{})
		for key, val := range extraMap {
			data[key] = val
		}
	}

	return data, nil
}

func (N Object) ToData() (interface{}, error) {
	data := map[string]interface{}{}

	for _, field := range N.Fields {
		fieldData, err := field.ToData()
		if err != nil {
			return nil, err
		}

		fieldMap := fieldData.(map[string]interface{})
		for key, val := range fieldMap {
			data[key] = val
		}
	}

	return data, nil
}

func (N Field) ToData() (interface{}, error) {
	val, err := N.Value.ToData()
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		N.Key.Value: val,
	}

	return data, nil
}

func (N Array) ToData() (interface{}, error) {
	data := []interface{}{}

	for _, elem := range N.Elems {
		elemData, err := elem.ToData()
		if err != nil {
			return nil, err
		}

		data = append(data, elemData)
	}

	return data, nil
}

func (N PathExpr) ToData() (interface{}, error) {
	data := map[string]interface{}{}

	return data, nil
}

func (N TokenPath) ToData() (interface{}, error) {
	data := map[string]interface{}{}

	return data, nil
}

func (N RangeExpr) ToData() (interface{}, error) {
	data := map[string]interface{}{}

	return data, nil
}

func (N BracePath) ToData() (interface{}, error) {
	data := map[string]interface{}{}

	return data, nil
}

func (N Token) ToData() (interface{}, error) {
	return N.Value, nil
}

func (N Integer) ToData() (interface{}, error) {
	return N.Value, nil
}

func (N Decimal) ToData() (interface{}, error) {
	return N.Value, nil
}

func (N Bool) ToData() (interface{}, error) {
	return N.Value, nil
}
