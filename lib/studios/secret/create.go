package secret

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/util"
)

const secretCreateQuery = `
mutation {
  secretCreateOneFor(values:{
    name:"{{name}}"
  }) {
    secretEverything {
      id
			createdAt
      name
    }
		message
		errors {
		  message
		}
  }
}
`

const secretCreateOutput = `
{{{data}}}
`

func Create(name, file string) error {

	vars := map[string]interface{}{
		"name": name,
	}

	data, err := util.SendRequest(secretCreateQuery, vars)
	if err != nil {
		return err
	}

	output, err := util.RenderString(secretCreateOutput, data)
	fmt.Println(output)

	err = Update(name, file)
	if err != nil {
		return err
	}

	return nil
}
