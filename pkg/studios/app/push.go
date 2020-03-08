package app

import (
	"errors"
	"fmt"

	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/pkg/config"
	"github.com/hofstadter-io/hof/pkg/util"
)

func Push() error {

	data, err := util.TarFiles(AppFiles, "./")
	if err != nil {
		fmt.Println("err", err)
		return err
	}

	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey
	host := util.ServerHost() + "/studios/app/push"
	acct, name := util.GetAcctAndName()

	req := gorequest.New().Post(host).
		Query("devmode=yes").
		Query("name="+name).
		Query("account="+acct).
		Set("apikey", apikey).
		Type("multipart").
		SendFile(data)

	resp, body, errs := req.End()
	// fmt.Println(resp, body, errs)

	if len(errs) != 0 || resp.StatusCode >= 500 {
		return errors.New("Internal Error: " + body)
	}
	if resp.StatusCode >= 400 {
		return errors.New("Bad Request: " + body)
	}

	fmt.Println(body)
	return nil
}
