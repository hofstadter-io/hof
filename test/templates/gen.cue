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
		{ a: val: 1 },
		{ b: val: 2 },
		{ c: val: 3 },
	]

	Map: {
		a: { val: 1 }
		b: { val: 2 }
		c: { val: 3 }
	}
}

#C: {
	foo: string
	i: int
	s: {
		b: bool | *false
		n: number & >0.0
	}
}

TestGen: #TestGen & {
	@gen(test)
  In: {
    Val: #A 
  }
	Val: #C
}


#TestGen: gen.#HofGenerator & {
  Outdir: "output"

  PackageName: ""

  In: {
    Val: _
    ...
  }

	Val: {...}

  Out: [...gen.#HofGeneratorFile] & [
    // Defaults
    {
      TemplateContent: "Val.a = '{{ .Val.a }}'\n"
      Filepath: "\(Outdir)/default.txt"
    },
    // Alternate delims
    {
      TemplateContent: "Val.a = '{% .Val.a %}'\n"
      Filepath: "\(Outdir)/altdelim.txt"
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
      Filepath: "\(Outdir)/named-things.txt"
    },

    // File based
    {
      TemplatePath: "template-file.txt"
      Filepath: "\(Outdir)/template-file.txt"
    },
    {
      TemplatePath: "template-altfile.txt"
      Filepath: "\(Outdir)/template-altfile.txt"
    },

    // User file
    {
      TemplateContent: "User file: '{{ file \"userfile.txt\" }}'\n"
      Filepath: "\(Outdir)/user-file.txt"
    },

		for F in _listFiles { F },
		for F in mapFiles { F },

		{
			DatafileFormat: "cue" 
			Filepath: "\(Outdir)/datafile-cue.cue"
		},
		{
			DatafileFormat: "json" 
			Filepath: "\(Outdir)/datafile-json.json"
			Val: #A
		},

  ]

	
	_listFiles: [ for idx, elem in In.Val.List {
			{
				TemplatePath: "template-elems.txt"
				Filepath: "\(Outdir)/list-elem-\(idx).txt"
				In: {
					Elem: elem
				}
			}
		}],

	mapFiles: [ for key, elem in In.Val.Map {
			{
				TemplatePath: "template-elems.txt"
				Filepath: "\(Outdir)/map-elem-\(key).txt"
				In: {
					Elem: elem
				}
			}
		}],

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
