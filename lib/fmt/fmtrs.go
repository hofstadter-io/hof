package fmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	gofmt "go/format"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue/format"
	"github.com/clbanning/mxj"
	"github.com/docker/docker/api/types"
	"github.com/naoina/toml"
	"gopkg.in/yaml.v3"

	"github.com/hofstadter-io/hof/cmd/hof/verinfo"
)

const ContainerPrefix = "hof-fmt-"

var defaultVersion = "dirty"

func init() {
	v := verinfo.Version
	if v != "Local" {
		defaultVersion = "v" + v
	}

	ov := os.Getenv("HOF_FMT_VERSION")
	if ov != "" {
		defaultVersion = ov
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
	// prettier
	".js":      "prettier/babel",
	".jsx":     "prettier/babel",
	".ts":      "prettier/typescript",
	".tsx":     "prettier/typescript",
	".graphql": "prettier/graphql",
	".yml":     "prettier/yaml",
	".html":    "prettier/html",
	".css":     "prettier/css",
	".less":    "prettier/less",
	".scss":    "prettier/scss",
	".md":      "prettier/markdown",
	".vue":     "prettier/vue",

	// black
	".py": "black/py",
}

var fmtrDefaultConfigs = map[string]interface{}{
	"prettier/babel": map[string]interface{}{
		"semi": false,
		"parser": "babel",
	},
	"prettier/typescript": map[string]interface{}{
		"parser": "typescript",
	},
	"prettier/graphql": map[string]interface{}{
		"parser": "graphql",
	},
	"prettier/html": map[string]interface{}{
		"parser": "html",
	},
	"prettier/css": map[string]interface{}{
		"parser": "css",
	},
	"prettier/less": map[string]interface{}{
		"parser": "less",
	},
	"prettier/scss": map[string]interface{}{
		"parser": "scss",
	},
	"prettier/markdown": map[string]interface{}{
		"parser": "markdown",
	},
	"prettier/vue": map[string]interface{}{
		"parser": "vue",
	},
	"prettier/handlebars": map[string]interface{}{
		"parser": "glimmer",
	},
	"prettier/go-template": map[string]interface{}{
		"parser": "go-template",
	},

	"black/py": map[string]interface{}{
		"parser": "go-template",
	},
}

func FormatSource(filename string, content []byte, fmtrName string, config interface{}, formatData bool) ([]byte, error) {
	// extract filename & extension
	// TODO, better extract multipart extensions (as supported by prettier)
	_, fn := filepath.Split(filename)
	ext := filepath.Ext(fn)

	// short-circuit builtin mime-types
	switch ext {
		case ".go":
			return gofmt.Source(content)

		case ".cue":
			if formatData {
				return formatCue(content)
			} else {
				return content, nil
			}

		case ".json":
			if formatData {
				return formatJson(content)
			} else {
				return content, nil
			}

		case ".yml", ".yaml":
			if formatData {
				return formatYaml(content)
			} else {
				return content, nil
			}

		case ".xml":
			if formatData {
				return formatXml(content)
			} else {
				return content, nil
			}

		case ".toml":
			if formatData {
				return formatToml(content)
			} else {
				return content, nil
			}

	}

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
	if resp.StatusCode >= 400 {
		fmt.Println("\n" + string(body) + "\n")
		return content, fmt.Errorf("error while formatting %s", filename)
	}

	content = body

	return content, nil
}

func formatCue(input []byte) ([]byte, error) {
	bs, err := format.Source(input)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func formatJson(input []byte) ([]byte, error) {
	v := make(map[string]interface{})
	err := json.Unmarshal(input, &v)
	if err != nil {
		return nil, err
	}	

	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func formatYaml(input []byte) ([]byte, error) {
	v := make(map[string]interface{})
	err := yaml.Unmarshal(input, &v)
	if err != nil {
		return nil, err
	}	

	bs, err := yaml.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func formatToml(input []byte) ([]byte, error) {
	v := make(map[string]interface{})
	err := toml.Unmarshal(input, &v)
	if err != nil {
		return nil, err
	}

	bs, err := toml.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func formatXml(input []byte) ([]byte, error) {
	xmlReader := bytes.NewReader(input)
	mv, err := mxj.NewMapXmlReader(xmlReader)
	if err != nil {
		return nil, err
	}

	bs, err := mv.XmlIndent("", "  ")
	if err != nil {
		return nil, err
	}
	return bs, nil
}
