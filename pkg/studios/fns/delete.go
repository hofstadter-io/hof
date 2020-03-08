package fns

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/util"
)

const fnDeleteQuery = `
mutation {
  fnDeleteOneFor(id:"{{id}}") {
    fnEverything {
      name
      id
      version
      type
			createdAt
    }
		message
		errors {
		  message
		}
  }
}
`

const fnDeleteOutput = `
{{{data}}}
`

func Delete(input string) error {

	var data interface{}
	var err error

	if util.IsValidUUID(input) {
		data, err = DeleteById(input)
	} else {
		data, err = DeleteByName(input)
	}

	if err != nil {
		return err
	}

	output, err := util.RenderString(fnDeleteOutput, data)

	fmt.Println(output)
	return err
}

func DeleteById(id string) (interface{}, error) {
	fmt.Println("DeleteById:", id)
	vars := map[string]interface{}{
		"id": id,
	}

	return util.SendRequest(fnDeleteQuery, vars)
}

func DeleteByName(name string) (interface{}, error) {
	fmt.Println("DeleteByName:", name)
	res, err := FilterByName(name)
	if err != nil {
		return nil, err
	}
	fmt.Println("Result:", res)

	basePath := "data.fnGetManyFor.fnEverything"

	id, err := util.FindIdFromName(basePath, name, fnListOutput, res)
	if err != nil {
		return nil, err
	}
	fmt.Println("ID:", id)

	return DeleteById(id)
}
