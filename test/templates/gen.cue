package templates

import (
  "github.com/hofstadter-io/hof/schema"
)

A :: {
  a: "a"
  N: {
    x: "x"
    y: "y"
  }
}

HofGenTest: TestGen & { In: Val: A }

TestGen :: schema.HofGenerator & {

  In: {
    Val: _
    ...
  }

  Out: [
    // Defaults
    schema.HofGeneratorFile & {
      Template: "Val.a = '{{ .Val.a }}'\n"
      Filepath: "default.txt"
    },
    // Alternate delims
    schema.HofGeneratorFile & {
      Template: "Val.a = '{% .Val.a %}'\n"
      Filepath: "altdelim.txt"
      AltDelims: true
      LHS2_D: "{%"
      RHS2_D: "%}"
      LHS3_D: "{%%"
      RHS3_D: "%%}"
    },
    // Swap delims, using defaults delims for swap/temp
    schema.HofGeneratorFile & {
      Template: "Val.a = '{% .Val.a %}' and also this should stay {{ .Hello }}\n"
      Filepath: "swapdelim.txt"
      AltDelims: true
      SwapDelims: true
      LHS2_D: "{%"
      RHS2_D: "%}"
      LHS3_D: "{%%"
      RHS3_D: "%%}"
    },
    // TODO Swap delims, using custom delims for swap/temp

    // Mustache system
    schema.HofGeneratorFile & {
      Template: "Val.a = '{{ Val.a }}'\n"
      TemplateSystem: "mustache"
      Filepath: "mustache.txt"
    },


    // Named things
    schema.HofGeneratorFile & {
      TemplateName: "named"
      Filepath: "named-things.txt"
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
    "static.txt": """
    Hello, I am a static file
    """
  }
}
