package fns

import (
	"errors"
	"fmt"

	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/lib/config"
	"github.com/hofstadter-io/hof/lib/util"
)

func Push() error {
	data, err := util.TarFiles(FuncFiles, "./")
	if err != nil {
		fmt.Println("err", err)
		return err
	}

	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey
	host := util.ServerHost() + "/studios/fns/push"
	acct, fname := util.GetAcctAndName()

	fmt.Println("Pushing:", fname)

	req := gorequest.New().Post(host).
		Query("account="+acct).
		Query("name="+fname).
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
