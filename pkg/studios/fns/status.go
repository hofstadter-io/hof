package fns

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/lib/config"
	"github.com/hofstadter-io/hof/lib/util"
)

func Status(fname string) error {
	if fname == "" {
		dir, _ := os.Getwd()
		fname = filepath.Base(dir)
	}

	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey
	host := util.ServerHost() + "/studios/fns/status"
	acct, _ := util.GetAcctAndName()

	resp, body, errs := gorequest.New().Get(host).
		Query("account="+acct).
		Query("name="+fname).
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
