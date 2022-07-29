package templates

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"text/template"

	"github.com/clbanning/mxj"
	"github.com/codemodus/kace"
	"github.com/hofstadter-io/hof/lib/dotpath"
	"github.com/kr/pretty"
	"github.com/naoina/toml"
	"gopkg.in/yaml.v3"
)

func AddGolangHelpers(t *template.Template) *template.Template {
	return t.Funcs(funcMap)
}

var funcMap = template.FuncMap{
	"json": Helper_json,
	"yaml": Helper_yaml,
	"toml": Helper_toml,
	"xml":  Helper_xml,

	"indent": Helper_indent,
	"pprint": Helper_pretty,
	"pretty": Helper_pretty,

	"lower":  Helper_lower,
	"upper":  Helper_upper,
	"title":  Helper_title,
	"pascal": Helper_pascal,
	"camel":  Helper_camel,
	"Camel":  Helper_camelT,
	"camelT": Helper_camelT,
	"snake":  Helper_snake,
	"SNAKE":  Helper_snakeU,
	"snakeU": Helper_snakeU,
	"kebab":  Helper_kebab,
	"KEBAB":  Helper_kebabU,
	"kebabU": Helper_kebabU,

	"concat":         Helper_concat,
	"join":           Helper_join,
	"contains":       Helper_contains,
	"split":          Helper_split,
	"replace":        Helper_replace,
	"hasprefix":      Helper_hasprefix,
	"hassuffix":      Helper_hassuffix,
	"trimspace":      Helper_trimspace,
	"trimprefix":     Helper_trimprefix,
	"trimsuffix":     Helper_trimsuffix,
	"trimto":         Helper_trimto_first,
	"trimfrom":       Helper_trimfrom_first,
	"trimto_first":   Helper_trimto_first,
	"trimfrom_first": Helper_trimfrom_first,
	"trimto_last":    Helper_trimto_last,
	"trimfrom_last":  Helper_trimfrom_last,
	"substr":         Helper_substr,
	"getprefix":      Helper_getprefix,
	"getsuffix":      Helper_getsuffix,
	"getbetween":     Helper_getbetween,

	"identity": Helper_identity,
	"dict":     Helper_dict,
	"file":     Helper_file,

	"gokind":  Helper_gokind,
	"builtin": Helper_builtin,

	"lookup": Helper_lookup,
	"dref":   Helper_lookup,
}

func Helper_yaml(value interface{}) string {
	bytes, err := yaml.Marshal(value)
	if err != nil {
		return err.Error()
	}
	ret := string(bytes)
	return ret
}

func Helper_toml(value interface{}) string {
	bytes, err := toml.Marshal(value)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

func Helper_json(value interface{}) string {
	bytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

// jsonl too?
func Helper_jsoninline(value interface{}) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

func Helper_xml(value interface{}) string {
	mv := mxj.Map(value.(map[string]interface{}))
	bytes, err := mv.XmlIndent("", "  ")
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

func Helper_concat(ss ...string) string {
	S := ""
	for _, s := range ss {
		S += s
	}
	return S
}

func Helper_join(sep string, ss ...string) string {
	return strings.Join(ss, sep)
}

func Helper_indent(value, indent string) string {
	ret := ""
	lines := strings.Split(value, "\n")
	for _, line := range lines {
		ret += indent + line + "\n"
	}
	ret = strings.TrimSuffix(ret, "\n")
	return ret
}

func Helper_pretty(value interface{}) string {
	return fmt.Sprintf("%# v", pretty.Formatter(value))
}

func Helper_lower(value string) string {
	return strings.ToLower(value)
}

func Helper_upper(value string) string {
	return strings.ToUpper(value)
}

func Helper_title(value string) string {
	return strings.Title(value)
}

func Helper_pascal(value string) string {
	return kace.Pascal(value)
}

func Helper_camel(value string) string {
	return kace.Camel(value)
}

func Helper_camelT(value string) string {
	return strings.Title(kace.Camel(value))
}

func Helper_snake(value string) string {
	return kace.Snake(value)
}

func Helper_snakeU(value string) string {
	return kace.SnakeUpper(value)
}

func Helper_kebab(value string) string {
	return kace.Kebab(value)
}

func Helper_kebabU(value string) string {
	return kace.KebabUpper(value)
}

func Helper_contains(str, srch string) string {
	if strings.Contains(str, srch) {
		return "true"
	}
	return ""
}

func Helper_split(str, sep string) []string {
	return strings.Split(str, sep)
}

func Helper_replace(str, old, new string, cnt int) string {
	return strings.Replace(str, old, new, cnt)
}

func Helper_hasprefix(str, pre string) string {
	if strings.HasPrefix(str, pre) {
		return "true"
	}
	return ""
}

func Helper_hassuffix(str, suf string) string {
	if strings.HasSuffix(str, suf) {
		return "true"
	}
	return ""
}

func Helper_trimspace(str string) string {
	return strings.TrimSpace(str)
}

func Helper_trimprefix(str, pre string) string {
	return strings.TrimPrefix(str, pre)
}

func Helper_trimsuffix(str, suf string) string {
	return strings.TrimSuffix(str, suf)
}

func Helper_trimto_first(str, pre string, keep bool) string {
	pos := strings.Index(str, pre)
	if pos >= 0 {
		if keep {
			return str[pos:]
		}
		return str[pos+len(pre):]
	}
	return str
}

func Helper_trimfrom_first(str, pre string, keep bool) string {
	pos := strings.Index(str, pre)
	if pos >= 0 {
		if keep {
			return str[:pos+len(pre)]
		}
		return str[:pos]
	}
	return str
}

func Helper_trimto_last(str, pre string, keep bool) string {
	pos := strings.LastIndex(str, pre)
	if pos >= 0 {
		if keep {
			return str[pos:]
		}
		return str[pos+len(pre):]
	}
	return str
}

func Helper_trimfrom_last(str, pre string, keep bool) string {
	pos := strings.LastIndex(str, pre)
	if pos >= 0 {
		if keep {
			return str[:pos+len(pre)]
		}
		return str[:pos]
	}
	return str
}

func Helper_substr(str string, start, end int) string {
	if end == -1 {
		end = len(str)
	}
	return str[start:end]
}

func Helper_getprefix(str, suf string) string {
	pos := strings.Index(str, suf)
	if pos >= 0 {
		return str[:pos]
	}
	return str
}

func Helper_getsuffix(str, suf string) string {
	pos := strings.Index(str, suf)
	if pos >= 0 {
		return str[pos+1:]
	}
	return str
}

func Helper_getbetween(str, lhs, rhs string) string {
	lpos := strings.Index(str, lhs)
	rpos := strings.LastIndex(str, rhs)
	if lpos < 0 {
		lpos = 0
	} else {
		lpos += 1
	}
	if rpos < 0 {
		rpos = len(str)
	}
	return str[lpos:rpos]
}

func Helper_identity(thing interface{}) interface{} {
	return thing
}

func Helper_dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i+=2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func Helper_file(filename string) string {
	body, err := ioutil.ReadFile(filename)

	// return content when non-error
	if err == nil {
		return string(body)
	}

	// return error on fallthrough, as a string for intemplate prcessing
	return fmt.Sprintf("ERROR: %v", err)
}

func Helper_gokind(input interface{}) string {
	v := reflect.ValueOf(input)
	k := v.Kind()
	return k.String()
}

var known_builtins = map[string]struct{}{
	"bool":        struct{}{},
	"byte":        struct{}{},
	"error":       struct{}{},
	"float":       struct{}{},
	"float32":     struct{}{},
	"float64":     struct{}{},
	"complex64":   struct{}{},
	"complex128":  struct{}{},
	"int":         struct{}{},
	"int8":        struct{}{},
	"int16":       struct{}{},
	"int32":       struct{}{},
	"int64":       struct{}{},
	"uint":        struct{}{},
	"uint8":       struct{}{},
	"uint16":      struct{}{},
	"uint32":      struct{}{},
	"uint64":      struct{}{},
	"rune":        struct{}{},
	"string":      struct{}{},
	"object":      struct{}{},
	"interface{}": struct{}{},
}

func Helper_builtin(str string) interface{} {
	_, ok := known_builtins[str]
	if ok {
		return true
	}
	return nil
}

func Helper_lookup(path string, data interface{}) interface{} {
	if data == nil {
		return fmt.Sprint("Nil data supplied for " + path)
	}

	// if OpenAPI format, convert to dotpath
	if strings.HasPrefix(path, "#/") {
		path = path[2:]
		path = strings.Replace(path, "/", ".", -2)
	}

	obj, err := dotpath.Get(path, data, true)
	if err != nil {
		return fmt.Sprint("Error during path search: " + err.Error())
	}

	if obj == nil {
		return fmt.Sprint("Path not found: " + path + fmt.Sprintf("\n%+v", data))
	}

	return obj
}
