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

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"
	cuejson "cuelang.org/go/pkg/encoding/json"
	"github.com/clbanning/mxj"
	"github.com/docker/docker/api/types"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"

	"github.com/hofstadter-io/hof/cmd/hof/verinfo"
)

const ContainerPrefix = "hof-fmt-"

var defaultVersion = "dirty"
var FORMAT_DISABLED = false
var DOCKER_FORMAT_DISABLED = false


func init() {
	v := verinfo.Version
	if v != "Local" {
		defaultVersion = "v" + v
	}

	ov := os.Getenv("HOF_FMT_VERSION")
	if ov != "" {
		defaultVersion = ov
	}

	formatters = make(map[string]*Formatter)
	for _,fmtr := range fmtrNames {
		formatters[fmtr] = &Formatter{Name: fmtr, Version: defaultVersion}
	}

	val := os.Getenv("HOF_FORMAT_DISABLED")
	if val == "true" || val == "1" {
		FORMAT_DISABLED=true
		DOCKER_FORMAT_DISABLED=true
	}
	
	// gracefully init images / containers
	err := GracefulInit()
	if err != nil {
		DOCKER_FORMAT_DISABLED=true
	}

}

func GracefulInit() error {
	err := initDockerCli()
	if err != nil {
		return err
	}

	err = updateFormatterStatus()
	if err != nil {
		return err
	}

	return nil
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

var fmtrNames = []string{
	"black",
	"csharpier",
	"prettier",
}

// Map file extensions to formatters
var extToFmtr = map[string]string {
	// python
	".py": "black/py",

	// csharp
	".cs": "csharpier/cs",

	// prettier
	".js":      "prettier/babel",
	".jsx":     "prettier/babel",
	".ts":      "prettier/typescript",
	".tsx":     "prettier/typescript",
	".graphql": "prettier/graphql",
	".html":    "prettier/html",
	".css":     "prettier/css",
	".less":    "prettier/less",
	".scss":    "prettier/scss",
	".md":      "prettier/markdown",
	".vue":     "prettier/vue",

	// TODO, there are a lot of others built in as well

	// prettier plugins,
	// TODO probably a separate image?
	".java":    "prettier/java",
	".rb":      "prettier/ruby",
	".rs":      "prettier/rust",
	".php":     "prettier/php",

	// This one is buggy
	// ".groovy":  "prettier/groovy",
	// going to introduce a dedicated container with a groovy script / server backed by CodeNarc
}

// Map wellknown filenames to formatters
var filenameToFmtr = map[string]string {
	"Jenkinsfile": "prettier/groovy",
}

var fmtrDefaultConfigs = map[string]interface{}{
	// python only
	"black/py": map[string]interface{}{
		"parser": "",
	},

	// csharp only
	"csharpier/cs": map[string]interface{}{
		"parser": "",
	},

	// pretty common
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

	// pretty extra
	"prettier/java": map[string]interface{}{
		"parser": "java",
	},
	"prettier/groovy": map[string]interface{}{
		"parser": "groovy",
	},
	"prettier/ruby": map[string]interface{}{
		"parser": "ruby",
	},
	"prettier/rust": map[string]interface{}{
		"parser": "jinx-rust",
	},
	"prettier/php": map[string]interface{}{
		"parser": "php",
	},
}

func FormatSource(filename string, content []byte, fmtrName string, config interface{}, formatData bool) ([]byte, error) {
	// short circuit here, so we don't everywhere this function is used
	if FORMAT_DISABLED {
		return content, nil
	}

	// extract filename & extension
	_, fn := filepath.Split(filename)
	fileParts := strings.Split(fn, ".")
	// filename without extension
	fileBase := fileParts[0]
	fileExt := strings.Join(fileParts[1:], ".")
	ext := filepath.Ext(fn)

	// short-circuit builtin mime-types
	// what about when there are multiple dots? How do we handle like below? do we need to?
	// examples: foo.tf.json & foo.tmpl.go
	switch ext {
		case ".go":
			return gofmt.Source(content)

		case ".cue":
			return formatCue(content)
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

	// short circuit here, so we don't everywhere this function is used
	if DOCKER_FORMAT_DISABLED {
		return content, nil
	}

	fmtrTool := ""
	ok := false

	// formatter was manually set
	if fmtrName != "" {
		config, ok = fmtrDefaultConfigs[fmtrName]
		if !ok {
			return content, fmt.Errorf("unknown formatter %q", fmtrName)
		}

		parts := strings.Split(fmtrName, "/")
		fmtrTool = parts[0]
	} else {
		// infer the formatter from filename/ext
		fmtrPath, ok := "", false

		if fileExt == "" {
			// look for wellknown filenames
			fmtrPath, ok = filenameToFmtr[fileBase]
			if !ok {
				// todo, ext not supported, alert (and not error?)
				return content, nil
			}
		} else {
			// look for extension to config
			// want to prefer the longest fileExt, then fallback to wellknown filebase
			// might need to upgrade to loop
			
			// try longest filepath
			fmtrPath, ok = extToFmtr[fileExt]
			// try golang exn
			if !ok {
				fmtrPath, ok = extToFmtr[ext]
			}
			// try wellknown filepath
			if !ok {
				fmtrPath, ok = filenameToFmtr[fileBase]
			}

			// if still not ok, return original content, we won't format
			if !ok {
				// todo, ext not supported, alert (and not error?)
				return content, nil
			}
		}

		config, ok = fmtrDefaultConfigs[fmtrPath]
		if !ok {
			panic(fmt.Sprint(fmtrPath, "not found in fmtrDefaultConfig"))
		}

		parts := strings.Split(fmtrPath, "/")
		fmtrTool = parts[0]
	}

	// we have a formatter picked out
	fmtr := formatters[fmtrTool]

	// start the formatter if not running
	if !fmtr.Running {
		err := startContainer(fmtrTool)
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

	if !bytes.HasSuffix(content, []byte{'\n'}) {
		content = append(content, '\n')
	}

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
	x, err := cuejson.Unmarshal(input)
	if err != nil {
		return nil, err
	}	

	var c cue.Context
	v := c.BuildExpr(x)

	s, err := cuejson.Marshal(v)
	if err != nil {
		return nil, err
	}

	s, err = cuejson.Indent([]byte(s), "", "    ")
	if err != nil {
		return nil, err
	}
	s += "\n"

	return []byte(s), nil
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

	buf := new(bytes.Buffer)
	err = toml.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
	bs = append(bs, '\n')
	return bs, nil
}
