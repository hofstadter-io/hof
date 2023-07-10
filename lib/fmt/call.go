package fmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hofstadter-io/hof/lib/container"
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

func (fmtr *Formatter) WaitForRunning(retry int, delay time.Duration) error {
	// fmt.Println("wait-running.0:", fmtr.Name, fmtr.Status, fmtr.Running, fmtr.Ready)
	// return if already running
	if fmtr.Running {
		return nil
	}

	for i := 0; i < retry; i++ {
		containers, err := container.GetContainers(ContainerPrefix + fmtr.Name)
		if err != nil {
			return err
		}

		for _, container := range containers {
			// extract name
			name := container.Names
			name = strings.TrimPrefix(name, ContainerPrefix)
			// fmt.Println("wait-running:", fmtr.Name, name, container.State)
			if name == fmtr.Name {
				fmtr.Status = container.State
				break
			}
		}

		if fmtr.Status == "running" {
			fmtr.Running = true
			err = updateFormatterStatus()
			if err != nil {
				return err
			}
			return nil
		}

		time.Sleep(delay)
	}

	return fmt.Errorf("formatter %s never started", fmtr.Name)
}

func (fmtr *Formatter) WaitForReady(retry int, delay time.Duration) error {
	// fmt.Println("wait-ready.0:", fmtr.Name, fmtr.Status, fmtr.Running, fmtr.Ready)

	// return if already ready
	if fmtr.Ready {
		return nil
	}

	// return error if not running
	if !fmtr.Running {
		return fmt.Errorf("formatter %s is not running", fmtr.Name)
	}

	// get ready check payload
	p, ok := fmtrReady[fmtr.Name]
	if !ok {
		fmt.Printf("warn: formatter %s does not have a ready config\n", fmtr.Name)
		return nil
	}

	payload := p.(map[string]any)

	for i := 0; i < retry; i++ {
		_, err := fmtr.Call("ready-check", []byte(payload["source"].(string)), payload["config"])
		// fmt.Println("wait-ready:", i, fmtr.Name, fmtr.Ready, err)
		// if no error, then ready
		if err == nil {
			fmtr.Ready = true
			return nil
		}
		time.Sleep(delay)
	}

	return fmt.Errorf("formatter %s is not ready", fmtr.Name)
}
