package crun

import (
	"errors"
	"fmt"

	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/lib/config"
	"github.com/hofstadter-io/hof/lib/util"
)

func Shutdown(name string) error {

	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey
	host := util.ServerHost() + "/studios/crun/shutdown"
	acct, fname := util.GetAcctAndName()
	if name == "" {
		name = fname
	}

	resp, body, errs := gorequest.New().Get(host).
		Query("account="+acct).
		Query("name="+name).
		Set("apikey", apikey).
		End()

	if len(errs) != 0 || resp.StatusCode >= 500 {
		return errors.New("Internal Error: " + body)
	}
	if resp.StatusCode >= 400 {
		return errors.New("Bad Request: " + body)
	}

	fmt.Println(body)
	return nil
}
