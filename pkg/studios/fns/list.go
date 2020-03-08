package fns

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/util"
)

const fnListQuery = `
query {
	fnGetManyFor(
    offset:{{after}}
    limit:{{limit}}
		{{#if filters}}
		filters: {
		  {{#if filters.search}}search:"{{filters.search}}"{{/if}}
		}
		{{/if}}
	) {
		fnEverything {
			id
			createdAt
			name
			version
			state
		}
		errors {
		  message
		}
  }
}
`

const fnListOutput = `
Name                    Version     State       ID
=======================================================================================
{{#each data.fnGetManyFor.fnEverything as |FN|}}
{{pw FN.name 24 ~}}
{{pw FN.version 12 ~}}
{{pw FN.state 12 ~}}
{{FN.id}}
{{/each}}
`

func List() error {
	vars := map[string]interface{}{
		"after": "0",
		"limit": "25",
	}

	data, err := util.SendRequest(fnListQuery, vars)
	if err != nil {
		return err
	}

	output, err := util.RenderString(fnListOutput, data)

	fmt.Println(output)
	return err
}

func FilterByName(name string) (interface{}, error) {
	vars := map[string]interface{}{
		"after": "0",
		"limit": "25",
		"filters": map[string]string{
			"search": name,
		},
	}

	return util.SendRequest(fnListQuery, vars)
}
