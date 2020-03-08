package crun

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/lib/config"
	"github.com/hofstadter-io/hof/lib/util"
)

func Pull() error {

	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey
	host := util.ServerHost() + "/studios/crun/pull"
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

	err := util.UntarFiles(CrunFiles, filepath.Join("funcs", name), bodyBytes)
	if err != nil {
		fmt.Println("err", err)
		return err
	}

	return nil
}
