package crun

import (
	"errors"
	"fmt"

	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/lib/config"
	"github.com/hofstadter-io/hof/lib/util"
)

func Push(name string) error {
	data, err := util.TarFiles(CrunFiles, "./")
	if err != nil {
		fmt.Println("err", err)
		return err
	}

	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey
	host := util.ServerHost() + "/studios/crun/push"
	acct, fname := util.GetAcctAndName()
	if name == "" {
		name = fname
	}

	fmt.Println("Pushing:", name)

	req := gorequest.New().Post(host).
		Query("account="+acct).
		Query("name="+name).
		Set("apikey", apikey).
		Type("multipart").
		SendFile(data)

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
