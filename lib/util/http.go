package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/parnurzeal/gorequest"

	"github.com/hofstadter-io/hof/lib/config"
)

func BuildRequest(path string) *gorequest.SuperAgent {
	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey

	url := ServerHost() + path
	acct, name := GetAcctAndName()

	req := gorequest.New().Get(url).
		Query("name="+name).
		Query("account="+acct).
		Set("apikey", apikey)

	return req
}

func ServerHost() string {
	return config.GetCurrentContext().Host
}

func GetAcctAndName() (string, string) {
	account := config.GetCurrentContext().Account

	dir, _ := os.Getwd()
	name := filepath.Base(dir)
	return account, name
}

func SimpleGet(path string) error {

	req := BuildRequest(path)
	resp, body, errs := req.End()

	check := "http2: server sent GOAWAY and closed the connection"
	if len(errs) != 0 && !strings.Contains(errs[0].Error(), check) {
		fmt.Println("errs:", errs)
		fmt.Println("resp:", resp)
		fmt.Println("body:", body)
		return errs[0]
	}

	if len(errs) != 0 || resp.StatusCode >= 500 {
		return errors.New("Internal Error: " + body)
	}
	if resp.StatusCode >= 400 {
		return errors.New("Bad Request: " + body)
	}

	fmt.Println(body)
	return nil

}
