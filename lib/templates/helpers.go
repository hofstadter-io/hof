package templates

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"text/template"

	"github.com/clbanning/mxj"
	"github.com/codemodus/kace"
	"github.com/kr/pretty"
	"github.com/naoina/toml"
	"gopkg.in/yaml.v3"

	"github.com/hofstadter-io/hof/lib/chat"
	"github.com/hofstadter-io/hof/lib/dotpath"
)

func (T *Template) AddGolangHelpers() {
	// traditional helpers
	T.T = T.T.Funcs(funcMap)

	// chat helpers

	chatMap := template.FuncMap{
		"chat": T.Helper_chat(),
		"gen": T.Helper_gen(),
		"render": T.Helper_render(),
	}

	T.T = T.T.Funcs(chatMap)
}

func hchat(msg string, args map[string]any) (string, error) {
	// fmt.Println("HCHAT", args)

	isOpenai := true
	model := "gpt-3.5-turbo"
	if m, ok := args["model"]; ok {
		model = m.(string)
	}
	if strings.HasPrefix(model, "chat-") || model == "bard" {
		isOpenai = false
	}
	switch model {
	case "bard":
		model = "chat-bison"
	case "gpt3", "gpt-3":
		model = "gpt-3.5-turbo"
	case "gpt4":
		model = "gpt-4"
	}

	P, ok := args["params"]
	if !ok || P == nil {
		P = make(map[string]any)
	}
	params := P.(map[string]any)
	// fmt.Println("PARAMS:", params)

	msgs := make([]chat.Message,0)
	exas := make([]chat.Example,0)

	msgs = append(msgs, chat.Message{
		Role: "user",
		Content: msg,
	})

	if isOpenai {
		resp, err := chat.OpenaiChat(model, msgs, params)
		if err != nil {
			return resp, err
		}
		return resp, nil
	} else {
		resp, err := chat.GoogleChat(model, msgs, exas, params)
		if err != nil {
			return resp, err
		}
		if b, ok := args["debug"]; ok && b.(bool) {
			fmt.Println("BARD:", resp)
		}
		return resp, nil
	}
}

// returns the full response object
func (T *Template) Helper_chat() func(string, ...map[string]any) any {

	return func(msg string, args ...map[string]any) any {
		if len(args) == 0 {
			args = append(args, make(map[string]any))
		}
		arrrrgs := args[0]
		curr := T.Buf.String()
		input := curr + msg
		body, err := hchat(input, arrrrgs)
		if err != nil {
			return body + "\n" + fmt.Sprint(err)
		}

		data := map[string]any{}
		err = json.Unmarshal([]byte(body), &data)
		if err != nil {
			return fmt.Sprintf("%s\n%s\n", body, err)
		}

		return data
	}
}

// returns just the message
func (T *Template) Helper_gen() any {

	return func(msg string, args ...map[string]any) string {
		if len(args) == 0 {
			args = append(args, make(map[string]any))
		}
		arrrrgs := args[0]
		curr := T.Buf.String()
		input := curr + msg
		body, err := hchat(input, arrrrgs)
		if err != nil {
			return body + "\n" + fmt.Sprint(err)
		}

		isOpenai := true
		model := "gpt-3.5-turbo"
		if m, ok := arrrrgs["model"]; ok {
			model = m.(string)
		}
		if strings.HasPrefix(model, "chat-") || model == "bard" {
			isOpenai = false
		}

		if isOpenai {
			resp, err := chat.OpenaiExtractContent(body)
			if err != nil {
				return resp + "\n" + fmt.Sprint(err)
			}
			return resp
		} else {
			resp, err := chat.GoogleExtractContent(body)
			if err != nil {
				return resp + "\n" + fmt.Sprint(err)
			}
			return resp
		}
	}
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

	"add":     Helper_add,
	"inc":     Helper_inc,

	"typeof":  Helper_gokind,
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

// Helper_indent indents lines of text, skipping the first line.
// The intended behavior is to match Helm's indent helper
func Helper_indent(indent interface{}, value string) string {
	ret := ""
	lines := strings.Split(value, "\n")

	// edge case
	if len(lines) == 1 {
		return value
	}
	
	// don't indent first line, left to user to place
	ret += lines[0] + "\n"
	lines = lines[1:]

	// indent, depending on arg
	switch i := indent.(type) {
		case string:
			for _, line := range lines {
				ret += i + line + "\n"
			}

		case int:
			spaces := strings.Repeat(" ", i)
			for _, line := range lines {
				ret += spaces + line + "\n"
			}

		default:
			return "indent only supports a string or integer as argument"
	}

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

	// return error on fallthrough, as a string for intemplate processing
	return fmt.Sprintf("ERROR: %v", err)
}

func Helper_add(lhs, rhs int) int {
	return lhs + rhs
}

func Helper_inc(val int) int {
	return val + 1
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

func Helper_builtin(str string) any {
	_, ok := known_builtins[str]
	if ok {
		return true
	}
	return nil
}

func Helper_lookup(path string, data any) any {
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

// todo, should we support turning the content back to an objecct?
// perhaps better to have different functions for this
func (T *Template) Helper_render() (func (name string, data any) any) {

	return func(name string, data any) any {
		t := T.T.Lookup(name)

		var b bytes.Buffer

		err := t.Execute(&b, data)
		if err != nil {
			return err.Error()
		}

		return b.String()
	}
}
