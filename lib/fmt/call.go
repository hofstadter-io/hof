package fmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// this file has functions for
// calling formatters in docker

func (fmtr *Formatter) Call(filename string, content []byte, config any) ([]byte, error) {
	if DOCKER_FORMAT_DISABLED {
		return content, nil
	}

	data := make(map[string]interface{})
	data["source"] = string(content)
	data["config"] = config

	// fmt.Println("fmtr.Call", fmtr.Name, fmtr.Port, data)

	bs, err := json.Marshal(data)
	if err != nil {
		return content, err
	}

	host := "http://localhost"
	if val := os.Getenv("HOF_FMT_HOST"); val != "" {
		host = val
	} else if fmtr.Host != "" {
		host = fmtr.Host
	}

	url := host + ":" + fmtr.Port

	if debug {
		fmt.Printf("fmt calling (%s) %s\n", fmtr.Name, url)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	if err != nil {
		return nil, fmt.Errorf("http new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("response Body:", string(body))
		return content, err
	}
	if resp.StatusCode >= 400 {
		fmt.Println("\n" + string(body) + "\n")
		return content, fmt.Errorf("error while formatting %s", filename)
	}

	content = body
	return content, nil
}

