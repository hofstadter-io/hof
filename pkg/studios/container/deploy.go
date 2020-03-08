package crun

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/lib/config"
	"github.com/hofstadter-io/hof/lib/util"
)

func Deploy(name string, push bool, memory string, concurrency int, timeout string, envs string) error {

	if push {
		err := Push(name)
		if err != nil {
			return err
		}
	}

	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey
	host := util.ServerHost() + "/studios/crun/deploy"
	acct, fname := util.GetAcctAndName()
	if name == "" {
		name = fname
	}

	fmt.Println("Building & Deploying...", name)

	req := gorequest.New().Post(host).
		Query("account=" + acct).
		Query("name=" + name).
		Query("memory=" + memory).
		Query("concurrency=" + fmt.Sprint(concurrency)).
		Query("timeout=" + timeout).
		Query("envs=" + envs)

	resp, body, errs := req.
		Timeout(20*time.Minute).
		Retry(0, 0*time.Second, http.StatusInternalServerError).
		Set("apikey", apikey).
		End()

	if len(errs) != 0 {
		fmt.Println(errs)
		return errors.New("Client Error")
	}
	if resp.StatusCode >= 500 {
		return errors.New("Internal Error: " + body)
	}
	if resp.StatusCode >= 400 {
		return errors.New("Bad Request: " + body)
	}

	fmt.Println(body)
	return nil
}
