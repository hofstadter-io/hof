package basic

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
    schema.HofGeneratorFile & {
      Template: "Val.a = '{{ .Val.a }}'\n"
      Filepath: "default.txt"
      AltDelims: false
    },
    schema.HofGeneratorFile & {
      Template: "Val.a = '{% .Val.a %}'\n"
      Filepath: "altdelim.txt"
      AltDelims: true
      LHS2_D: "{%"
      RHS2_D: "%}"
      LHS3_D: "{%%"
      RHS3_D: "%%}"
    },
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
  ]
}
