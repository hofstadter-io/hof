package util

import (
	"encoding/json"
	"fmt"

	"github.com/hofstadter-io/dotpath"
	"github.com/hofstadter-io/hof/lib/config"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

func SendRequest(queryTemplate string, vars interface{}) (interface{}, error) {
	ctx := config.GetCurrentContext()
	apikey := ctx.APIKey
	host := ServerHost() + "/graphql"
	acct, _ := GetAcctAndName()

	query, err := RenderString(queryTemplate, vars)
	if err != nil {
		return nil, err
	}

	send := map[string]interface{}{
		"query":     query,
		"variables": nil,
	}

	req := gorequest.New().Post(host).
		Query("account="+acct).
		Set("apikey", apikey).
		Send(send)

	resp, body, errs := req.EndBytes()

	if len(errs) != 0 || resp.StatusCode >= 500 {
		return nil, errors.New("Internal Error: " + string(body))
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New("Bad Request: " + string(body))
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func FindIdFromName(basePath, name, listOutput string, res interface{}) (string, error) {
	path := fmt.Sprintf("%s.[name==%s]", basePath, name)
	elem, err := dotpath.Get(path, res, false)
	if err != nil {
		return "", err
	}

	fmt.Println("Elem:", path, elem)

	if elem == nil || len(elem.([]interface{})) == 0 {
		fmt.Println("not found, see results:")
		output, err := RenderString(listOutput, res)
		if err != nil {
			return "", err
		}
		fmt.Println(output)
		fmt.Println("--- end results ---")
		return "", errors.New("not found")
	}

	path = fmt.Sprintf("%s.[name==%s].id", basePath, name)
	id, err := dotpath.Get(path, res, false)
	if err != nil {
		return "", err
	}

	ID, ok := id.(string)
	if !ok {
		return "", errors.New("ID Not String")
	}

	return ID, nil
}
