package app

import (
	"errors"
	"fmt"

	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/pkg/config"
	"github.com/hofstadter-io/hof/pkg/util"
)

func Pull() error {
	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey

	host := util.ServerHost() + "/studios/app/pull"
	acct, name := util.GetAcctAndName()

	resp, bodyBytes, errs := gorequest.New().Get(host).
		Query("name="+name).
		Query("account="+acct).
		Set("apikey", apikey).
		EndBytes()

	if len(errs) != 0 || resp.StatusCode >= 500 {
		return errors.New("Internal Error: " + fmt.Sprint(errs))
	}
	if resp.StatusCode >= 400 {
		return errors.New("Bad Request: " + fmt.Sprint(errs))
	}

	err := util.UntarFiles(AppFiles, ".", bodyBytes)
	if err != nil {
		fmt.Println("err", err)
		return err
	}

	fmt.Println("App pulled locally")
	return nil
}
