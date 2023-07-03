package fmt

import (
	"bytes"
	"fmt"
	gofmt "go/format"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	cuejson "cuelang.org/go/pkg/encoding/json"
	cueyaml "cuelang.org/go/pkg/encoding/yaml"
	"github.com/BurntSushi/toml"
	"github.com/clbanning/mxj"

	"github.com/hofstadter-io/hof/cmd/hof/verinfo"
	"github.com/hofstadter-io/hof/lib/container"
)

const ContainerPrefix = "hof-fmt-"

var (
	defaultVersion         = "dirty"
	FORMAT_DISABLED        = false
	DOCKER_FORMAT_DISABLED = false
)

// (CONSIDER) make this comma separated, so we can have fallback?
var CONTAINER_REPO = "ghcr.io/hofstadter-io"

var debug = false

func init() {
	v := verinfo.Version
	if v != "Local" {
		if !strings.HasPrefix(v, "v") {
			v = "v" + v
		}
		defaultVersion = v
	}

	ds := os.Getenv("HOF_FMT_DEBUG")
	if ds != "" {
		dv, err := strconv.ParseBool(ds)
		if err != nil {
			fmt.Println("Error parsing HOF_FMT_DEBUG:", err)
		} else {
			debug = dv
		}
	}

	ov := os.Getenv("HOF_FMT_VERSION")
	if ov != "" {
		defaultVersion = ov
	}

	formatters = make(map[string]*Formatter)
	for _, fmtr := range fmtrNames {
		formatters[fmtr] = &Formatter{Name: fmtr, Version: defaultVersion}
	}

	val := os.Getenv("HOF_FMT_DISABLED")
	if val != "" {
		dv, err := strconv.ParseBool(ds)
		if err != nil {
			fmt.Println("Error parsing HOF_FMT_DISABLED:", err)
		} else {
			FORMAT_DISABLED = dv
			DOCKER_FORMAT_DISABLED = dv
		}
	}

	hr := os.Getenv("HOF_FMT_REGISTRY")
	if hr != "" {
		CONTAINER_REPO = hr
	}

	// gracefully init images / containers
	err := Init()
	if err != nil {
		if debug {
			fmt.Println("fmt init error:", err)
		}
		DOCKER_FORMAT_DISABLED = true
	}

	if debug {
		fmt.Println("FORMAT_DISABLED", FORMAT_DISABLED)
		fmt.Println("DOCKER_FORMAT_DISABLED", DOCKER_FORMAT_DISABLED)
	}
}

func Init() error {
	err := container.InitClient()
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
	Name      string
	Version   string
	Available []string

	// Info
	Status    string
	Running   bool
	Ready     bool
	Host      string
	Port      string
	Container *container.Container
	Images    []*container.Image

	Config  interface{}
	Default interface{}
}

var formatters map[string]*Formatter

var fmtrNames = []string{
	"black",
	"csharpier",
	"prettier",
}

var fmtrEnvs = map[string][]string{
	"black":     nil,
	"csharpier": nil,
	"prettier": {
		"PRETTIER_RUBY_TIMEOUT_MS=10000",
	},
}

var fmtrReady = map[string]any{
	"black": map[string]any{
		"config": fmtrDefaultConfigs["black/py"],
		"source": "n = 1",
	},
	"csharpier": map[string]any{
		"config": fmtrDefaultConfigs["csharpier/cs"],
		"source": "var n = 1;",
	},
	"prettier": map[string]any{
		"config": fmtrDefaultConfigs["prettier/js"],
		"source": "var n = 1;",
	},
}

// Map file extensions to formatters
var extToFmtr = map[string]string{
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
	".java": "prettier/java",
	".php":  "prettier/php",
	".rb":   "prettier/ruby",
	".rs":   "prettier/rust",
	".sql":  "prettier/sql",

	// This one is buggy
	// ".groovy":  "prettier/groovy",
	// going to introduce a dedicated container with a groovy script / server backed by CodeNarc
}

// Map wellknown filenames to formatters
var filenameToFmtr = map[string]string{
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
		"semi":   false,
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
	"prettier/sql": map[string]interface{}{
		"parser": "sql",
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
		err := Start(fmtrTool, false)
		if err != nil {
			return content, err
		}
	}

	fmtd, err := fmtr.Call(filename, content, config)
	if err != nil {
		return content, err
	}
	content = fmtd

	// add a final newline if not present
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

	s, err = cuejson.Indent([]byte(s), "", "  ")
	if err != nil {
		return nil, err
	}
	s += "\n"

	return []byte(s), nil
}

func formatYaml(input []byte) ([]byte, error) {
	expr, err := cueyaml.Unmarshal(input)
	if err != nil {
		return nil, err
	}

	ctx := cuecontext.New()

	v := ctx.BuildExpr(expr)

	s, err := cueyaml.Marshal(v)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
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
