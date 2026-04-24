package veg


models: [n=string]: {
  @agentic(model)
  name: string| *n
  id: string | *n
  provider: string | *"vertex" | "openai"

  baseurl: string | *""

  // other settings
}

models: {
  "gemini-3-flash": id: "gemini-3-flash-preview"
  "gemini-3-pro": id: "gemini-3.1-pro-preview"

  "gemma-4": { id: "gemma4:31b", provider: "openai", baseurl: "http://nitrogen-lan:11434/v1"}
  "qwen-3.6": { id: "qwen3.6:35b-a3b", provider: "openai", baseurl: "http://nitrogen-lan:11434/v1"}

  // "gemini-2.5-flash-lite": id: "gemini-2.5-flash-lite-preview-09-2025"
  // "gemini-2.5-flash": id: "gemini-2.5-flash-preview-09-2025"
  // "gemini-2.5-pro": id: "gemini-2.5-pro"

  // third party need some extra registration (wonder if implementation too...?)
  // "kimi-k2-thinking": id: "moonshot/kimi-k2-thinking-maas"
  // "deepseek-r1": id: "deepseek-ai/deepseek-r1-0528-maas"
}

