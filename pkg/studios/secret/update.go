package secret

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/pkg/config"
	"github.com/hofstadter-io/hof/pkg/util"
)

func Update(name, file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("Error: file " + file + " does not exist")
		return nil
	}

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	fmt.Println("Updating Secret:", name, file)
	fmt.Println(string(contents))

	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey
	host := util.ServerHost() + "/studios/secrets/push"
	acct := config.GetCurrentContext().Account

	resp, body, errs := gorequest.New().Post(host).
		Query("account="+acct).
		Query("name="+name).
		Set("apikey", apikey).
		Type("text").
		Send(string(contents)).
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
