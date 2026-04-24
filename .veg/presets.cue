package veg

presets: [string]: {
  @agentic(preset)
  agent: string
  model: string
  env: string

  // you can also define envs that
  // always work from a remote git
  // just one dir
  dir: string
  // multiple dirs?
  // dirs: [...string]
}

presets: default: {
  agent: "code_assist"
  model: "qwen-3.6"
  env:   "veg-dev"
  dir:   "."
}