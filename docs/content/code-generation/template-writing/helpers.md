---
title: "Template Helpers"
weight: 10
---


## Helpers

In addition to Go's `text/template` defaults,
`hof` adds the following functions

There are a number of additional helpers provided.

Notes:

- Many of these are based on the Go stdlib functions with the same name.
- Args can be literal, a .Variable, or the result of another function. See Go's documentation
- It can be advantageous to do as many manipulations in CUE rather than in the template system.


`text/template` documentation can be found here:
https://pkg.go.dev/text/template



### Encode to a data format:

- `json` (`{{ json .Data }}`)
- `yaml` (`{{ yaml .Data }}`)
- `toml` (`{{ toml .Data }}`)
- `xml` (`{{ xml .Data }}`)

### Value helpers:

- `dict` (`{{ ... (dict "key1" "val1" "key2" .Value ...) }}`) - create a map on the fly, useful for passing multiple values to partials
- `lookup` (`{{ lookup "path.to.value" .Value }}`) - returns the value at [hofstadter-io/dotpath](https://github.com/hofstadter-io/dotpath/blob/master/examples/test.go) from the given value
- `pretty` (`{{ pretty .Value }}`) uses https://github.com/kr/pretty
- `printf` (`{{ printf "%s %d" .String .Number }}`)

### Casings conversion helpers:

- `lower` (`{{ lower .String }}`) (lowercase)
- `upper` (`{{ upper .String }}`) (UPPERCASE)
- `title` (`{{ title .String }}`) (Titlecasestring
- `pascal` (`{{ pascal .String }}`) (PascalCaseString)
- `camel` (`{{ camel .String }}`) (camelCaseString)
- `Camel` (`{{ Camel .String }}`) (camel + title)
- `snake` (`{{ snake .String }}`) (snake_case_string)
- `SNAKE` (`{{ SNAKE .String }}`) (SNAKE_CASE_STRING)
- `kebab` (`{{ kebab .String }}`) (kebab-case-string)
- `KEBAB` (`{{ KEBAB .String }}`) (KEBAB-CASE-STRING)

### String processing helpers:

- `concat` (`{{ concat "foo" "bar" }}`)
- `join` (`{{ join ", " "foo" "bar" }}`)
- `contains` (`{{ contains .String .Substr }}`)
- `split` (`{{ split .String .Separator }}`)
- `replace` (`{{ replace .String .Old .New .Count }}`)
- `hasprefix` (`{{ hasprefix .String .Prefix }}`)
- `hassuffix` (`{{ hassuffix .String .Suffix }}`)
- `trimprefix` (`{{ trimprefix .String .Prefix }}`)
- `trimsuffix` (`{{ trimsuffix .String .Suffix }}`)
- `trimto` (alias for trimto_first)
- `trimfrom` (alias for trimfrom_first)
- `trimto_first` (`{{ trimto_first .String .Substr .Keep }}`) Trim a string until the first Substring, possibly keeping the substring
- `trimfrom_first` (`{{ trimfrom_first .String .Substr .Keep  }}`)
- `trimto_last` (`{{ trimto_last .String .Substr .Keep }}`)
- `trimfrom_last` (`{{ trimfrom_last .String .Substr .Keep }}`)
- `substr` (`{{ substr .String .StartInt .EndInt }}`)
- `getprefix` (`{{ getprefix .String .Substr }}`) (opposite of `trimprefix`)
- `getsuffix` (`{{ getsuffix .String .Substr }}`)
- `getbetween` (`{{ getbetween .String .LhsSubstr .RhsSubstr }}`)

The `indent` helper is disucssed in the
[section on indentation](/code-generation/template-writing/indentation/).

### Other helpers:

- `identity` (`{{ identity .Value }}`) - returns Value, ah maths
- `gokind` (`{{ gokind .Value }}`) - returns the Go reflect.Kind as a string
- `builtin` (`{{ builtin }}`) - returns true if the string is a known Go builtin
- `file` (`{{ file "path/to/file.txt" }}`) - loads a file from disk

