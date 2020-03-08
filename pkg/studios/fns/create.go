package fns

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hofstadter-io/hof/lib/extern"
	"github.com/hofstadter-io/hof/lib/util"
)

const fnCreateQuery = `
mutation {
  fnCreateOneFor(values:{
    name:"{{name}}"
    version:"{{version}}"
    type:"{{type}}"
  }) {
    fnEverything {
      name
      id
      version
      type
			createdAt
    }
  }
}
`

const fnCreateOutput = `
{{{data}}}
`

func Create(name, template string, here bool) error {
	var version string
	var err error

	if name == "" {
		_, fname := util.GetAcctAndName()
		name = fname
	}

	if template != "none" {

		if template == "" || template[0] == '#' || template[0] == '@' {
			template = "https://github.com/hofstadter-io/studios-functions" + template
		}

		url, version, subpath := extern.SplitParts(template)

		data := map[string]interface{}{}
		data["FuncName"] = name

		var dir string

		if strings.HasPrefix(url, "https") {
			dir, err = util.CloneRepo(url, version)
			if err != nil {
				return err
			}
		} else {
			// assume local, just copy, so working copy
			dir = url
		}

		if here {
			err = util.RenderDirNameSub(filepath.Join(dir, subpath), ".", data)
		} else {
			err = util.RenderDirNameSub(filepath.Join(dir, subpath), name, data)
		}
		if err != nil {
			return err
		}

	}

	vars := map[string]interface{}{
		"name":    name,
		"type":    "starter",
		"version": version,
	}

	data, err := util.SendRequest(fnCreateQuery, vars)
	if err != nil {
		return err
	}

	output, err := util.RenderString(fnCreateOutput, data)

	fmt.Println(output)
	return err
}
