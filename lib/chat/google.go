package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/parnurzeal/gorequest"
)

var googleApiKey string
var googleProject string

func GoogleChat(model string, messages []Message, examples []Example, params map[string]any) (string, error) {
	if !(strings.HasPrefix(model, "chat-bison")) {
		return "", fmt.Errorf("incompatible google chat model: %q", model)
	}

	var err error

	if googleApiKey == "" {
		// load from gcloud
		googleApiKey, err = getGoogleConfig([]string{"gcloud", "auth", "print-access-token"})
		if err != nil {
			return "", fmt.Errorf("while acquiring gcloud auth: %w", err)
		}
	}
	if googleProject == "" {
		// load from gcloud
		googleProject, err = getGoogleConfig([]string{"gcloud", "config", "get", "project"})
		if err != nil {
			return "", fmt.Errorf("while acquiring gcloud auth: %w", err)
		}
	}

	inst := make(map[string]any)
	// process messages
	if messages[0].Role == "system" {
		inst["context"] = messages[0].Content
		messages = messages[1:]
	}
	msgs := []map[string]any{}
	for _, M := range messages {
		m := map[string]any{
			"author": M.Role,
			"content": M.Content,
		}
		msgs = append(msgs, m)
	}
	inst["messages"] = msgs

	// process examples
	ex := []map[string]any{}
	for _, E := range examples {
		ex = append(ex, map[string]any{
			"input": E.Input,
			"output": E.Output,
		})
	}
	inst["examples"] = ex

	// process parameters
	prms := make(map[string]any)
	if params != nil {
		prms["temperature"] = params["Temperature"]
		prms["maxOutputTokens"] = params["MaxTokens"]
		prms["topP"] = params["TopP"]
		prms["topK"] = params["TopK"]
	}

	data := make(map[string]any)
	data["instances"] = []map[string]any{inst}
	data["parameters"] = prms

	R := gorequest.New()
	R.Header.Add("Content-Type", "application/json")
	R.Header.Add("Authorization", "Bearer " + googleApiKey)

	url := "https://us-central1-aiplatform.googleapis.com/v1/projects/%s/locations/us-central1/publishers/google/models/%s:predict"
	R.Url = fmt.Sprintf(url, googleProject, model)
	R.Method = "POST"

	// fmt.Printf("%#+v\n", data)
	R.Send(data)

	resp, body, errs := R.End()

	if len(errs) != 0 && resp == nil {
		return body, fmt.Errorf("%v", errs)
	}
	if len(errs) != 0 {
		return body, fmt.Errorf("Internal Error:\n%v\n%s\n", errs, body)
	}

	return body, nil
}

func getGoogleConfig(args []string) (string, error){
	cmd := exec.Command(args[0], args[1:]...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func GoogleExtractContent(body string) (string, error) {

	d := map[string]any{}
	err := json.Unmarshal([]byte(body), &d)
	if err != nil {
		return body, err
	}

	p, ok := d["predictions"]
	if !ok {
		return body, fmt.Errorf("Error: no predictions found")
	}
	P, ok := p.([]any)
	if !ok {
		return body, fmt.Errorf("Error: predictions in not a list")
	}
	P0, ok := P[0].(map[string]any)
	if !ok {
		return body, fmt.Errorf("Error: prediction 0 not a map")
	}

	c, ok := P0["candidates"]
	if !ok {
		return body, fmt.Errorf("Error: no candidates found")
	}
	C, ok := c.([]any)
	if !ok {
		return body, fmt.Errorf("Error: canidates in not a list")
	}
	C0, ok := C[0].(map[string]any)
	if !ok {
		return body, fmt.Errorf("Error: candidate 0 not a map")
	}

	content, ok := C0["content"]
	if !ok {
		return body, fmt.Errorf("Error: no content found")
	}

	ret, ok := content.(string)
	if !ok {
		return body, fmt.Errorf("Error: content is not a string")
	}

	return ret, nil
}
