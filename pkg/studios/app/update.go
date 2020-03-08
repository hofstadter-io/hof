package app

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"

	"github.com/hofstadter-io/hof/lib/config"
	"github.com/hofstadter-io/hof/lib/util"
)

func Update(version string) error {
	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey
	host := util.ServerHost() + "/studios/app/update"

	resp, body, errs := gorequest.New().Get(host).
		Query("version="+version).
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
