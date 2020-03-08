package crun

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/util"
)

const crunListQuery = `
query {
	crunGetManyFor(
    offset:{{after}}
    limit:{{limit}}
		{{#if filters}}
		filters: {
		  {{#if filters.search}}search:"{{filters.search}}"{{/if}}
		}
		{{/if}}
	) {
		crunEverything {
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

const crunListOutput = `
Name                    Version     State       ID
=======================================================================================
{{#each data.crunGetManyFor.crunEverything as |CRUN|}}
{{pw CRUN.name 24 ~}}
{{pw CRUN.version 12 ~}}
{{pw CRUN.state 12 ~}}
{{CRUN.id}}
{{/each}}
`

func List() error {
	vars := map[string]interface{}{
		"after": "0",
		"limit": "25",
	}
	data, err := util.SendRequest(crunListQuery, vars)
	if err != nil {
		return err
	}

	output, err := util.RenderString(crunListOutput, data)

	fmt.Println(output)
	return err
}

func FilterByName(name string) (interface{}, error) {
	vars := map[string]interface{}{
		"after": "0",
		"limit": "25",
		"filters": map[string]string{
			"name": name,
		},
	}

	return util.SendRequest(crunListQuery, vars)
}
