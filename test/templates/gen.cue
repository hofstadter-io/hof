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

TestGen: #TestGen & {
	@gen(test)
  In: {
    Val: #A 
  }
}


#TestGen: gen.#HofGenerator & {
  Outdir: "output"

  PackageName: ""

  In: {
    Val: _
    ...
  }

  Out: [...gen.#HofGeneratorFile] & [
    // Defaults
    {
      TemplateContent: "Val.a = '{{ .Val.a }}'\n"
      Filepath: "\(Outdir)/default.txt"
      TemplateConfig: {
        Engine: "golang"
      }
    },
    // Alternate delims
    {
      TemplateContent: "Val.a = '{% .Val.a %}'\n"
      Filepath: "\(Outdir)/altdelim.txt"
      TemplateConfig: {
				AltDelims: {
					LHS2: "{%"
					RHS2: "%}"
					LHS3: "{%%"
					RHS3: "%%}"
				}
      }
    },
    // Swap delims, using defaults delims for swap/temp
    {
      TemplateContent: "Val.a = '{% .Val.a %}' and also this should stay {{ .Hello }}\n"
      Filepath: "\(Outdir)/swapdelim.txt"
      TemplateConfig: {
        TmpDelims: true
				AltDelims: {
					LHS2: "{%"
					RHS2: "%}"
					LHS3: "{%%"
					RHS3: "%%}"
				}
      }
    },
    // TODO Swap delims, using custom delims for swap/temp

    // Mustache system
    {
      TemplateContent: "Val.a = '{{ Val.a }}'\n"
      Filepath: "\(Outdir)/mustache.txt"
      TemplateConfig: {
        Engine: "raymond"
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
	}, {
		Globs: ["templates/altdelim-*"]
		Config: {
      TmpDelims: true
			AltDelims: {
				LHS2: "{%"
				RHS2: "%}"
				LHS3: "{%%"
				RHS3: "%%}"
			}
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
