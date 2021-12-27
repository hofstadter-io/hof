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
	CUE: #C
}


#TestGen: gen.#HofGenerator & {
  Outdir: "output"

  PackageName: ""

  In: {
    Val: _
    ...
  }

	CUE: {...}

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

		// TODO
		// Per-template In, also in repeated
		// Repeated Files
  ]

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
