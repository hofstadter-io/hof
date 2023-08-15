package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/cmd/hof/verinfo"
	hfmt "github.com/hofstadter-io/hof/lib/fmt"
	"github.com/hofstadter-io/hof/lib/templates"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func InitModule(name string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) error {
	module := "hof.io"
	if strings.Contains(name,"/") {
		i := strings.LastIndex(name,"/")
		module, name = name[:i], name[i+1:]
	}
	// possibly extract explicit package
	pkg := name
	if strings.Contains(name,":") {
		i := strings.LastIndex(name,":")
		name, pkg = name[:i], name[i+1:]
	}
	fmt.Printf("Initializing: %s/%s in pkg %s\n", module, name, pkg)

	ver := verinfo.HofVersion
	if !strings.HasPrefix(ver, "v") {
		ver = "v" + ver
	}

	// construct template input data
	data := map[string]interface{}{
		"CueVer": verinfo.CueVersion,
		"HofVer": ver,
		"Module": module,
		"Name": name,
		"Package": pkg,
	}

	// local helper to render and write embedded templates
	render := func(outpath, content string) error {
		if rootflags.Verbosity > 0 {
			fmt.Println("rendering:", outpath)
		}
		ft, err := templates.CreateFromString(outpath, content, templates.Delims{})
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
			bs, err = hfmt.FormatSource(outpath, bs, "", nil, true)
			if err != nil {
				return err
			}
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

	err := render(name + ".cue", newModuleTopTemplate)
	if err != nil {
		return err
	}
	err = render("gen/gen.cue", newModuleGenTemplate)
	if err != nil {
		return err
	}
	err = render("cue.mod/module.cue", cuemodFileTemplate)
	if err != nil {
		return err
	}
	// todo, fetch deps
	msg, err := yagu.Shell("hof mod tidy", "")
	fmt.Println(msg)
	if err != nil {
		return err
	}
	// make some dirs
	dirs := []string{"templates", "partials", "statics", "examples", "creators", "gen", "schema"}
	for _, dir := range dirs {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	err = render("-", finalMsg)
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) adhocAsModule() error {
	name := R.GenFlags.AsModule
	module := "hof.io"
	if strings.Contains(name,"/") {
		i := strings.LastIndex(name,"/")
		module, name = name[:i], name[i+1:]
	}
	// possibly extract explicit package
	pkg := name
	if strings.Contains(name,":") {
		i := strings.LastIndex(name,":")
		name, pkg = name[:i], name[i+1:]
	}
	fmt.Printf("Initializing: %s/%s in pkg %s\n", module, name, pkg)

	// parse template flags
	tcfgs  := []AdhocTemplateConfig{}
	tfiles := make([]string,0)
	for _, tf := range R.GenFlags.Template {
		cfg, err := ParseTemplateFlag(tf)
		if err != nil {
			return err
		}
		tcfgs  = append(tcfgs, cfg)
		tfiles = append(tfiles, cfg.Filepath)

		if R.Flags.Verbosity > 0 {
			fmt.Printf("%#v\n", cfg)
		}
	}

	// top-level fields that would have been accessible
	ins := []string{}

	// the function that decides if a top-level field
	// should be added to the generator templates
	// this is necessary so that the refrenced inputs still work
	keep := func(value cue.Value) bool {
		path := value.Path().String()
		for _, node := range R.Nodes {
			// if we find a match, we need to decide
			if path == node.Hof.Path {
				// exceptions to the rule below
				// - it is a datamodel
				if node.Hof.Datamodel.Root {
					return true
				}

				// don't keep anything with a $hof field
				return false
			}
		}
		return true
	}

	// get top-level CUE value as a struct
	S, err := R.Value.Struct()
	if err != nil {
		return err
	}

	// Loop through all top level fields
	// They must be regular by design
	iter := S.Fields()
	for iter.Next() {
		if keep(iter.Value()) {
			ins = append(ins, iter.Label())
		}
	}

	// get generator names that were loaded by -G
	gens := []string{}
	for _, G := range R.Generators {
		label := G.Hof.Label
		if label == "AdhocGen" {
			// this is probably no longer the case
			continue
		}
		gens = append(gens, label)
	}
	sort.Strings(gens)

	ver := verinfo.HofVersion
	if !strings.HasPrefix(ver, "v") {
		ver = "v" + ver
	}

	// construct template input data
	data := map[string]interface{}{
		"Configs": tcfgs,
		"CueVer": verinfo.CueVersion,
		"Diff3": R.GenFlags.Diff3,
		"Entrypoints": R.Entrypoints,
		"Generators": gens,
		"HofVer": ver,
		"Inputs": ins,
		"Module": module,
		"Name": name,
		"Outdir": R.GenFlags.Outdir,
		"Package": pkg,
		"Partials": R.GenFlags.Partial,
		"Templates": tfiles,
		"WatchFast": R.GenFlags.WatchFast,
		"WatchFull": R.GenFlags.WatchFull,
	}

	// local helper to render and write embedded templates
	render := func(outpath, content string) error {
		if R.Flags.Verbosity > 0 {
			fmt.Println("rendering:", outpath)
		}
		ft, err := templates.CreateFromString(outpath, content, templates.Delims{})
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
			bs, err = hfmt.FormatSource(outpath, bs, "", nil, true)
			if err != nil {
				return err
			}

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

	if R.Flags.Verbosity > 0 {
		fmt.Println("writing:", name)
	}
	if name == "-" {
		err = render(name, asModuleTopTemplate)
		if err != nil {
			return err
		}
		err = render(name, asModuleGenTemplate)
		if err != nil {
			return err
		}
	} else {
		err = render(name + ".cue", asModuleTopTemplate)
		if err != nil {
			return err
		}
		err = render("gen/gen.cue", asModuleGenTemplate)
		if err != nil {
			return err
		}
		err = render("cue.mod/module.cue", cuemodFileTemplate)
		if err != nil {
			return err
		}

		// fetch deps
		cmd := exec.Command("hof", "mod", "tidy")
		out, err := cmd.CombinedOutput()
		fmt.Println(string(out))
		if err != nil {
			return err
		}

		err = render("-", finalMsg)
		if err != nil {
			return err
		}
	}

	return nil
}

const asModuleTopTemplate = `
package {{ snake .Package }}

import (
	"{{ .Module }}/{{ .Name }}/gen"
)

// This is example usage of your generator
{{ camelT .Name }}: gen.Generator & {
	@gen({{ .Name }})

	// inputs to the generator
	{{ range .Inputs -}}
	"{{.}}": {{.}},
	{{ else }}
	// you almost certainly need to
	// manually add input data here
	// Data: data
	{{- end }}

	// other settings
	Diff3: {{ .Diff3 }}
	{{ if .Outdir }}
	Outdir: "{{ .Outdir }}"
	{{ end }}
	
	{{ if .WatchFull }}
	// File globs to watch and trigger regen when changed
	// Normally, a user would set this to their designs / datamodel
	WatchFull: [ {{ range .WatchFull }}"{{.}}", {{ end }} ]
	{{ end }}
	{{ if .WatchFast }}
	// This is really only useful for module authors
	WatchFast:  [ {{ range .WatchFast  }}"{{.}}", {{ end }} ]
	{{ end }}

	// required by examples inside the same module
	// your users do not set or see this field
	ModuleName: ""
}
`


const asModuleGenTemplate = `
package gen

import (
	"github.com/hofstadter-io/hof/schema/gen"
)

// This is your reusable generator module
Generator: gen.Generator & {

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
	// Data: _
	// Input: #Input
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
	ModuleName: string | *"{{ .Module }}/{{ .Name }}"

	{{ if .Templates -}}
	// Templates: [gen.#Templates & {Globs: ["./templates/**/*"], TrimPrefix: "./templates/"}]
	Templates: [ { Globs: [ {{ range .Templates }}"{{.}}", {{ end }} ] } ]
	{{ else }}
	// Templates: [gen.#Templates & {Globs: ["./templates/**/*"], TrimPrefix: "./templates/"}]
	Templates: []
	{{ end }}
	{{ if .Partials -}}
	// Partials: [gen.#Templates & {Globs: ["./partials/**/*"], TrimPrefix: "./partials/"}]
	Partials:  [ { Globs: [ {{ range .Partials  }}"{{.}}", {{ end }} ] } ]
	{{ else }}
	// Partials: [gen.#Templates & {Globs: ["./partials/**/*"], TrimPrefix: "./partials/"}]
	Partials: []
	{{ end }}
	Statics: []

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


const initMsg = `To run the '{{.Name}}' generator...
  $ hof gen        ... or ...
  $ hof gen{{range .Entrypoints}} {{.}}{{ end }} {{ .Name }}.cue -G {{ .Name }}
`
const newModuleTopTemplate = `
package {{ snake .Package }}

import (
	"{{ .Module }}/{{ .Name }}/gen"
)

// This is example usage of your generator
{{ camelT .Name }}: gen.Generator & {
	@gen({{ .Name }})

	// inputs to the generator
	Data: { ... }
	Outdir: "./out/"
	
	// File globs to watch and trigger regen when changed
	// Normally, a user would set this to their designs / datamodel
	WatchFull: [...string]
	// This is helpful when authoring generator modules
	WatchFast:  [...string]

	// required by examples inside the same module
	// your users do not set or see this field
	ModuleName: ""
}
`

const newModuleGenTemplate = `
package gen

import (
	"github.com/hofstadter-io/hof/schema/gen"
)

// This is your reusable generator module
Generator: gen.Generator & {

	//
	// user input fields
	//

	// this is the interface for this generator module
	// typically you enforce schema(s) here
	// Data: _
	// Input: #Input

	//
	// Internal Fields
	//

	// This is the global input data the templates will see
	// You can reshape and transform the user inputs
	// While we put it under internal, you can expose In
	// or you can omit In and skip having a global context
	In: {
		// fill as needed
		...
	}

	// required for hof CUE modules to work
	// your users do not set or see this field
	ModuleName: string | *"{{ .Module }}/{{ .Name }}"

	// The final list of files for hof to generate
	// fill this with file values
	Out: [...gen.#File] & [
	]

	// you can create any intermediate values you need internally

	// open, so your users can build on this
	...
}
`

const cuemodFileTemplate = `
module: "{{ .Module }}/{{ .Name }}"
cue: "{{ .CueVer }}"

require: {
	"github.com/hofstadter-io/hof": "{{ .HofVer }}"
}
`

const finalMsg = `To run the '{{.Name}}' generator...
  $ hof gen        ... or ...
  $ hof gen{{range .Entrypoints}} {{.}}{{ end }} {{ .Name }}.cue -G {{ .Name }}
`

