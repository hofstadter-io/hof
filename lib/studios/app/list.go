package app

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/util"
)

const appListQuery = `
query {
	appGetManyFor(
    offset:{{after}}
    limit:{{limit}}
		{{#if filters}}
		filters: {
		  {{#if filters.search}}search:"{{filters.search}}"{{/if}}
		}
		{{/if}}
	) {
		appStatus {
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

const appListOutput = `
Name                    Version     State       ID
=======================================================================================
{{#each data.appGetManyFor.appStatus as |APP|}}
{{pw APP.name 24 ~}}
{{pw APP.version 12 ~}}
{{pw APP.state 12 ~}}
{{APP.id}}
{{/each}}
`

func List() error {

	data, err := GetList()
	if err != nil {
		return err
	}

	output, err := util.RenderString(appListOutput, data)

	fmt.Println(output)
	return err
}

func GetList() (interface{}, error) {
	vars := map[string]interface{}{
		"after": "0",
		"limit": "25",
	}

	return util.SendRequest(appListQuery, vars)
}

func FilterByName(name string) (interface{}, error) {
	vars := map[string]interface{}{
		"after": "0",
		"limit": "25",
		"filters": map[string]string{
			"name": name,
		},
	}

	return util.SendRequest(appListQuery, vars)
}
