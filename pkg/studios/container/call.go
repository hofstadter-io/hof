package crun

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/lib/config"
	"github.com/hofstadter-io/hof/lib/util"
)

func Call(name string, data string) error {

	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey
	host := util.ServerHost() + "/studios/crun/call"
	acct, _ := util.GetAcctAndName()

	// fmt.Println("Calling:", host)

	req := gorequest.New().Post(host).
		Query("account="+acct).
		Query("name="+name).
		Set("apikey", apikey)

	if len(data) > 0 {
		if data[:1] == "@" {
			bytes, err := ioutil.ReadFile(data[1:])
			if err != nil {
				return err
			}
			data = string(bytes)
		}

		req.Send(data)
	}

	resp, body, errs := req.End()

	if len(errs) != 0 || resp.StatusCode >= 500 {
		return errors.New("Internal Error: " + body)
	}
	if resp.StatusCode >= 400 {
		return errors.New("Bad Request: " + body)
	}

	fmt.Println(body)
	return nil
}
