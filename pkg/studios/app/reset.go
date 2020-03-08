package app

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/config"
	"github.com/hofstadter-io/hof/lib/util"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

func Reset(name string) error {
	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey
	host := util.ServerHost() + "/studios/app/reset"

	acct, appname := util.GetAcctAndName()
	if name == "" {
		name = appname
	}

	req := gorequest.New().Get(host).
		Query("name="+name).
		Query("account="+acct).
		Set("apikey", apikey)

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
