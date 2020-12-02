package diff

import (
	"fmt"

	"github.com/pkg/errors"
)

func Diff(original interface{}, current interface{}) (diff interface{}, err error) {

	fmt.Println("DIFF'n types: " + fmt.Sprintf("%T, %T", original, current))
	// check that they are the same type at the root
	// If different - error
	// If same - go recurse
	switch original.(type) {

	case map[string]interface{}:
		_, ok := current.(map[string]interface{})
		if !ok {
			return nil, errors.New("undiffable types, not the same type" + fmt.Sprintf("%T, %T", original, current))
		}

		return rdiff(original, current)

	case []interface{}:
		_, ok := current.([]interface{})
		if !ok {
			return nil, errors.New("undiffable types, not the same type" + fmt.Sprintf("%T, %T", original, current))
		}

		return rdiff(original, current)

	default:
		// TODO check for golang types with reflect
		return nil, errors.New("undiffable original, must be map or slice" + fmt.Sprintf("%T, %+v", original, original))

	}

	return nil, errors.New("undiffable original" + fmt.Sprintf("%T, %+v", original, original))
}

func rdiff(original interface{}, current interface{}) (diff interface{}, err error) {

	switch O := original.(type) {

	case map[string]interface{}:
		C, ok := current.(map[string]interface{})
		if !ok {
			return nil, errors.New("undiffable types, not the same type" + fmt.Sprintf("%T, %T", original, current))
		}

		fmt.Println("diffable types" + fmt.Sprintf("%T, %T", O, C))

		return nil, nil


	case []interface{}:
		C, ok := current.([]interface{})
		if !ok {
			return nil, errors.New("undiffable types, not the same type" + fmt.Sprintf("%T, %T", original, current))
		}

		fmt.Println("diffable types" + fmt.Sprintf("%T, %T", O, C))

		// if elements have names,

		return nil, nil

	default:
		// TODO check for golang types with reflect
		return nil, errors.New("undiffable original, must be map or slice" + fmt.Sprintf("%T, %+v", original, original))

	}

	return nil, errors.New("undiffable known" + fmt.Sprintf("%#+v, %#+v", original, current))
}
