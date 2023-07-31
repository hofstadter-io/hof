package templates

import (
	"github.com/hofstadter-io/hof/schema/gen"
)

#A: {
	a: "a"
	b: "b"
	N: {
		x: "x"
		y: "y"
	}
	List: [
		{a: val: 1},
		{b: val: 2},
		{c: val: 3},
	]

	Map: {
		a: {val: 1}
		b: {val: 2}
		c: {val: 3}
	}
}

#C: {
	foo: string
	i:   int
	s: {
		b: bool | *false
		n: number & >0.0
	}
}

TestGen: #TestGen & {
	@gen(test)
	In: {
		Val:   #A
		maybe: bool | *false @tag(maybe,type=bool)
	}
	Val: #C
}

#TestGen: gen.#HofGenerator & {
	Outdir: "output"

	PackageName: ""

	In: {
		Val:   _
		maybe: bool
		...
	}

	Val: {...}

	Out: [...gen.#HofGeneratorFile] & [
		// Defaults
		{
			TemplateContent: "Val.a = '{{ .Val.a }}'\n"
			Filepath:        "default.txt"
		},
		// Alternate delims
		{
			TemplateContent: "Val.a = '{% .Val.a %}'\n"
			Filepath:        "altdelim.txt"
			TemplateDelims: {
				LHS: "{%"
				RHS: "%}"
			}
			In: {
				extra: "foobar"
			}
		},

		// Named things
		{
			TemplatePath: "named"
			Filepath:     "named-things.txt"
		},

		// File based
		{
			TemplatePath: "template-file.txt"
			Filepath:     "template-file.txt"
		},
		{
			TemplatePath: "template-altfile.txt"
			Filepath:     "template-altfile.txt"
		},

		// User file
		{
			TemplateContent: "User file: '{{ file \"userfile.txt\" }}'\n"
			Filepath:        "user-file.txt"
		},

		for F in _listFiles {F},
		for F in mapFiles {F},

		{
			DatafileFormat: "cue"
			Filepath:       "datafile-cue.cue"
		},
		{
			DatafileFormat: "json"
			Filepath:       "datafile-json.json"
			Val:            #A
		},

	]

	_listFiles: [ for idx, elem in In.Val.List {
		{
			TemplatePath: "template-elems.txt"
			Filepath:     "list-elem-\(idx).txt"
			In: {
				Elem: elem
			}
		}
	}]

	mapFiles: [ for key, elem in In.Val.Map {
		{
			TemplatePath: "template-elems.txt"
			Filepath:     "map-elem-\(key).txt"
			In: {
				Elem: elem
			}
		}
	}]

	Templates: [{
		Globs: ["templates/template-*"]
		TrimPrefix: "templates/"
	}, {
		Globs: ["./templates/altdelim-*"]
		TrimPrefix: "./templates/"
		Delims: {
			LHS: "{%"
			RHS: "%}"
		}
	}]

	EmbeddedTemplates: {
		named: {
			Content: """
				embedded template is '{{ .Val.a }}'
				"""
		}
	}

	EmbeddedPartials: {
		named: {
			Content: """
				embedded partial is '{{ .Val.a }}'
				"""
		}
	}

	EmbeddedStatics: {
		"static-cue.txt": """
			Hello, I am a static file in cue
			"""
	}
}
