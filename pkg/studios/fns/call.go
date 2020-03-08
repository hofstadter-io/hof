package fns

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/lib/config"
	"github.com/hofstadter-io/hof/lib/util"
)

func Call(fname string, data string) error {
	if fname == "" {
		dir, _ := os.Getwd()
		fname = filepath.Base(dir)
	}

	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey
	host := util.ServerHost() + "/studios/fns/call"
	acct, _ := util.GetAcctAndName()

	if data[:1] == "@" {
		bytes, err := ioutil.ReadFile(data[1:])
		if err != nil {
			return err
		}
		data = string(bytes)
	}

	req := gorequest.New().Post(host).
		Query("account="+acct).
		Query("name="+fname).
		Set("apikey", apikey).
		Send(data)

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
