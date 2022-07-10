package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/templates"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func (R *Runtime) AsModule() error {
	FP := R.Flagpole
	name := FP.AsModule

	if R.Verbosity > 0 {
		fmt.Println("modularizing", name)
		fmt.Println(strings.Repeat("-", 23))
	}

	// parse template flags
	tcfgs  := []AdhocTemplateConfig{}
	tfiles := make([]string,0)
	for _, tf := range R.Flagpole.Template {
		cfg, err := parseTemplateFlag(tf)
		if err != nil {
			return err
		}
		tcfgs  = append(tcfgs, cfg)
		tfiles = append(tfiles, cfg.Filepath)

		if R.Verbosity > 0 {
			fmt.Printf("%#v\n", cfg)
		}
	}

	// top-level fields that would have been accessible
	ins := []string{}
	// but not anything with this attribute
	filters := map[string]bool{
		"gen": true,
		"hof": true,
	}

	// get top-level CUE value as a struct
	S, err := R.CueRuntime.CueValue.Struct()
	if err != nil {
		return err
	}

	// Loop through all top level fields
	// They must be regular by design
	iter := S.Fields()
	for iter.Next() {

		// what we will add if not filtered
		label := iter.Label()

		// let's possibly filster
		value := iter.Value()
		attrs := value.Attributes(cue.ValueAttr)

		filtered := false
		// find top-level with gen attr
		for _, A := range attrs {
			// does it have "@gen()"
			if _, ok := filters[A.Name()]; ok {
				filtered = true
			}
		}

		if !filtered {
			ins = append(ins, label)
		}
	}

	// get generator names that were loaded by -G
	gens := []string{}
	for label := range R.Generators {
		if label == "AdhocGen" {
			continue
		}
		gens = append(gens, label)
	}
	sort.Strings(gens)

	// construct template input data
	data := map[string]interface{}{
		"Name": name,
		"Inputs": ins,
		"Configs": tcfgs,
		"Templates": tfiles,
		"Partials": FP.Partial,
		"Generators": gens,
		"Diff3": FP.Diff3,
		"WatchGlobs": FP.WatchGlobs,
		"WatchXcue": FP.WatchXcue,
	}

	// local helper to render and write embedded templates
	render := func(outpath, content string) error {
		if R.Verbosity > 0 {
			fmt.Println("rendering:", outpath)
		}
		ft, err := templates.CreateFromString(outpath, content, nil)
		if err != nil {
			return err
		}
		bs, err := ft.Render(data)
		if err != nil {
			return err
		}
		if outpath == "-" {
			fmt.Println(string(bs))
			return nil
		} else {
			if strings.Contains(outpath, "/") {
				dir, _ := filepath.Split(outpath)
				err := os.MkdirAll(dir, 0755)
				if err != nil {
					return err
				}
			}
			return os.WriteFile(outpath, bs, 0644)
		}
	}

	if R.Verbosity > 0 {
		fmt.Println("writing:", name)
	}
	if name == "-" {
		err = render(name, asModuleTemplate)
		if err != nil {
			return err
		}
	} else {
		err = render(name + ".cue", asModuleTemplate)
		if err != nil {
			return err
		}
		err = render("cue.mods", cuemodsTemplate)
		if err != nil {
			return err
		}
		err = render("cue.mod/module.cue", cuemodFileTemplate)
		if err != nil {
			return err
		}
		// todo, fetch deps
		msg, err := yagu.Bash("hof mod vendor cue")
		fmt.Println(msg)
		if err != nil {
			return err
		}

		err = render("-", finalMsg)
		if err != nil {
			return err
		}
	}

	if R.Verbosity > 0 {
		fmt.Println(strings.Repeat("-", 23))
	}

	return nil
}

const asModuleTemplate = `
package {{ .Name }}

import (
	"github.com/hofstadter-io/hof/schema/gen"
)

// This is example usage of your generator
{{ camelT .Name }}Example: #{{ camelT .Name }}Generator & {
	@gen({{ .Name }})

	// inputs to the generator
	{{ range .Inputs -}}
	"{{.}}": {{.}},
	{{ else }}
	// you almost certainly need to
	// manually add input data here
	// "data": data,
	{{- end }}

	// other settings
	Diff3: {{ .Diff3 }}	
	Outdir: "./"
	
	{{ if .WatchGlobs }}
	// File globs to watch and trigger regen when changed
	// Normally, a user would set this to their designs / datamodel
	WatchGlobs: [ {{ range .WatchGlobs }}"{{.}}", {{ end }} ]
	{{ end }}
	{{ if .WatchXcue }}
	// This is really only useful for module authors
	WatchXcue:  [...string] | *[ {{ range .WatchXcue  }}"{{.}}", {{ end }} ]
	{{ end }}

	// required by examples inside the same module
	// your users do not set or see this field
	PackageName: ""
}


// This is your reusable generator module
//
#{{ camelT .Name }}Generator: gen.#Generator & {

	//
	// user input fields
	//

	// this is the interface for this generator module
	// typically you enforce schema(s) here
	{{ range .Inputs -}}
	{{.}}: _
	{{ else }}
	// you almost certainly need to
	// manually add input fields & schemas here
	// "data": _
	// "input": #Input
	{{- end }}

	//
	// Internal Fields
	//

	// This is the global input data the templates will see
	// You can reshape and transform the user inputs
	// While we put it under internal, you can expose In
	In: {
		// if you want to user your input data
		// add top-level fields from your
		// CUE entrypoints here, adjusting as needed
		// Since you made this a module for others,
		// it won't output until this field is filled

		{{ range .Inputs -}}
		"{{.}}": {{ . }}
		{{ else }}
		// you almost certainly need to
		// manually add input fields & schemas here
		// "data": _
		// "input": #Input
		{{- end }}

		...
	}

	// required for hof CUE modules to work
	// your users do not set or see this field
	PackageName: string | *"hof.io/{{ .Name }}"

	{{ if .Templates -}}
	// Templates: [{Globs: ["./templates/**/*"], TrimPrefix: "./templates/"}]
	Templates: [ { Globs: [ {{ range .Templates }}"{{.}}", {{ end }} ] } ]
	{{ else }}
	// Templates: [{Globs: ["./templates/**/*"], TrimPrefix: "./templates/"}]
	Templates: []
	{{ end }}
	{{ if .Partials -}}
	// Partials: [#Templates & {Globs: ["./partials/**/*"], TrimPrefix: "./partials/"}]
	Partials:  [ { Globs: [ {{ range .Partials  }}"{{.}}", {{ end }} ] } ]
	{{ else }}
	// Partials: [#Templates & {Globs: ["./partials/**/*"], TrimPrefix: "./partials/"}]
	Partials: []
	{{ end }}

	{{ if .Generators -}}
	// these should be in the same CUE package
	// or you may have to manually import as needed
	Generators: {
		{{ range .Generators -}}
		"{{.}}": {{.}},
		{{- end }}
	}
	{{ end }}

	// The final list of files for hof to generate
	Out: [...gen.#File] & [
		{{ range $i, $cfg := .Configs -}}
		{{ if $cfg.Repeated }}for _, t in t_{{ $i }} { t }{{ else }}t_{{ $i }}{{end}},
		{{ end }}
	]

	// These are the -T mappings
	{{ range $i, $cfg := .Configs -}}
	t_{{ $i }}: {{ if not .Repeated }}{
		{{ if .Cuepath }}In: In.{{.Cuepath}}{{ end }}
		{{ if .Schema  }}In: {{.Schema}}{{end}}
		TemplatePath: "{{ .Filepath }}"
		Filepath:     "{{ .Outpath }}"
	}{{ else }}[ for _,el in In.{{.Cuepath}} {
		{{ if .Cuepath }}In: el{{ end }}
		{{ if .Schema  }}In: {{.Schema}}{{end}}
		TemplatePath: "{{ .Filepath }}"
		Filepath:     "{{ .Outpath }}"
	}]{{ end }}
	{{ end }}

	// so your users can build on this
	...
}
`

const cuemodFileTemplate = `
module: "hof.io/{{ .Name }}"
`

const cuemodsTemplate = `
module hof.io/{{ .Name }}

cue v0.4.3

require (
	github.com/hofstadter-io/hof v0.6.3
)
`

const finalMsg = `Now run hof gen -G {{ .Name }}`
