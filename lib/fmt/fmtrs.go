package fmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"

	"github.com/hofstadter-io/hof/cmd/hof/verinfo"
)

const ContainerPrefix = "hof-fmt-"

var defaultVersion = "dirty"

func init() {
	v := verinfo.Version
	if v != "Local" {
		defaultVersion = v
	}
}

func GracefulInit() {
	err := initDockerCli()
	if err != nil {
		return
	}

	err = updateFormatterStatus()
	if err != nil {
		return
	}
}

type Formatter struct {
	// name, same as tools/%
	Name    string
	Version string
	Available []string

	// Info
	Running   bool
	Port      string
	Container *types.Container
	Images    []*types.ImageSummary

	Config  interface{}
	Default interface{}
}

var formatters map[string]*Formatter

func init() {
	formatters = make(map[string]*Formatter)
	for _,fmtr := range fmtrNames {
		formatters[fmtr] = &Formatter{Name: fmtr, Version: defaultVersion}
	}
}

var fmtrNames = []string{
	"black",
	"prettier",
}

var extToFmtr = map[string]string {
	".js": "prettier/js",
}

var fmtrDefaultConfigs = map[string]interface{}{
	"prettier/js": map[string]interface{}{
		"semi": false,
		"parser": "babel",
	},
}

func FormatSource(filename string, content []byte, fmtrName string, config interface{}) ([]byte, error) {
	// extract filename & extension
	// TODO, better extract multipart extensions (as supported by prettier)
	_, fn := filepath.Split(filename)
	ext := filepath.Ext(fn)

	// if the users hadn't 
	if config == nil {
		// look for extension to config
		fmtrPath, ok := extToFmtr[ext]
		if !ok {
			// todo, ext not supported, alert (and not error?)
			return content, nil
		}

		config, ok = fmtrDefaultConfigs[fmtrPath]
		if !ok {
			panic(fmt.Sprint(fmtrPath, "not found in fmtrDefaultConfig"))
		}

		parts := strings.Split(fmtrPath, "/")
		fmtrName = parts[0]
	}

	// we have a formatter picked out
	fmtr := formatters[fmtrName]

	// start the formatter if not running
	if fmtr.Port == "" {
		err := startContainer(fmtrName)
		if err != nil {
			return content, err
		}
		err = updateFormatterStatus()
		if err != nil {
			return content, err
		}
	}

	switch fmtrName {
		case "prettier":
			data := make(map[string]interface{})
			data["source"] = string(content)
			data["config"] = config

			bs, err := json.Marshal(data)
			if err != nil {
				return content, err
			}

			url := "http://localhost:" + fmtr.Port

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
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

			content = body

			
		default: 
			panic(fmt.Sprint("unknown formatter", fmtrName))
	}

	return content, nil
}


