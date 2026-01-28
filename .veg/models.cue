package veg


models: [n=string]: {
  @agentic(model)
  name: string| *n
  id: string | *n

  // other settings
}

models: {
  "gemini-3-flash": id: "gemini-3-flash-preview"
  "gemini-3-pro": id: "gemini-3-pro-preview"

  // "gemini-2.5-flash-lite": id: "gemini-2.5-flash-lite-preview-09-2025"
  // "gemini-2.5-flash": id: "gemini-2.5-flash-preview-09-2025"
  // "gemini-2.5-pro": id: "gemini-2.5-pro"

  // third party need some extra registration (wonder if implementation too...?)
  // "kimi-k2-thinking": id: "moonshot/kimi-k2-thinking-maas"
  // "deepseek-r1": id: "deepseek-ai/deepseek-r1-0528-maas"
}
