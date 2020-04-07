package templates

import (
  "github.com/hofstadter-io/hof/schema"
)

A :: {
  a: "a"
  b: "b"
  N: {
    x: "x"
    y: "y"
  }
}

HofGenTest: TestGen & {
  In: {
    Val: A 
  }
  ...
}


TestGen :: schema.HofGenerator & {
  Outdir: "output"

  In: {
    Val: _
    ...
  }

  Out: [
    // Defaults
    schema.HofGeneratorFile & {
      Template: "Val.a = '{{ .Val.a }}'\n"
      Filepath: "\(Outdir)/default.txt"
      TemplateConfig: {
        TemplateSystem: "golang"
      }
    },
    // Alternate delims
    schema.HofGeneratorFile & {
      Template: "Val.a = '{% .Val.a %}'\n"
      Filepath: "\(Outdir)/altdelim.txt"
      TemplateConfig: {
        AltDelims: true
        LHS2_D: "{%"
        RHS2_D: "%}"
        LHS3_D: "{%%"
        RHS3_D: "%%}"
      }
    },
    // Swap delims, using defaults delims for swap/temp
    schema.HofGeneratorFile & {
      Template: "Val.a = '{% .Val.a %}' and also this should stay {{ .Hello }}\n"
      Filepath: "\(Outdir)/swapdelim.txt"
      TemplateConfig: {
        AltDelims: true
        SwapDelims: true
        LHS2_D: "{%"
        RHS2_D: "%}"
        LHS3_D: "{%%"
        RHS3_D: "%%}"
      }
    },
    // TODO Swap delims, using custom delims for swap/temp

    // Mustache system
    schema.HofGeneratorFile & {
      Template: "Val.a = '{{ Val.a }}'\n"
      Filepath: "\(Outdir)/mustache.txt"
      TemplateConfig: {
        TemplateSystem: "raymond"
      }
    },


    // Named things
    schema.HofGeneratorFile & {
      TemplateName: "named"
      Filepath: "\(Outdir)/named-things.txt"
    },

    // File based
    schema.HofGeneratorFile & {
      TemplateName: "template-file.txt"
      Filepath: "\(Outdir)/template-file.txt"
    },

    // User file
    schema.HofGeneratorFile & {
      Template: "User file: '{{ file \"userfile.txt\" }}'\n"
      Filepath: "\(Outdir)/user-file.txt"
    },
  ]

  TemplatesDir: "templates/"
  PartialsDir: "partials/"
  StaticGlobs: ["static/**"]

  NamedTemplates: {
    named: """
    named is '{{ .Val.a }}'
    """
  }

  NamedPartials: {
    named: """
    partial is '{{ .Val.a }}'
    """
  }

  StaticFiles: {
    "static-cue.txt": """
    Hello, I am a static file in cue
    """
  }
}
