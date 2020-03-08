package app

import (
	"fmt"
	"strings"

	"github.com/hofstadter-io/hof/lib/util"
)

const appCreateQuery = `
mutation {
  appCreateOneFor(values:{
    name:"{{name}}"
    version:"{{version}}"
    type:"{{type}}"
  }) {
    appEverything {
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

const appCreateOutput = `
{{{data}}}
`

func Create(name, template string, here bool) error {

	if !here {
		err := CreateLocal(name, template)
		if err != nil {
			return err
		}
	}

	return CreateRemote(name)
}

func CreateRemote(name string) error {
	// TODO read hof.yaml

	vars := map[string]interface{}{
		"name": name,
		"type": "starter",
	}

	// REMOTE CREATION
	retdata, err := util.SendRequest(appCreateQuery, vars)
	if err != nil {
		return err
	}

	output, err := util.RenderString(appCreateOutput, retdata)

	fmt.Println(output)
	return err
}

func CreateLocal(name, template string) error {
	var url, version, dir string
	var err error

	parts := strings.Split(template, "@")
	if len(parts) == 2 {
		url = parts[0]
		version = parts[1]
	} else {
		url = template
	}

	data := map[string]interface{}{}
	data["AppName"] = name

	if strings.HasPrefix(url, "https") {
		dir, err = util.CloneRepo(url, version)
		if err != nil {
			return err
		}
	} else {
		// assume local, just copy, so working copy
		dir = url
	}

	return util.RenderDirNameSub(dir, name, data)
}
