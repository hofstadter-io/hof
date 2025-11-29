package agents

agents: [n=string]: {
  name: string | *n
  model: string | *"gemini-2.5-flash"
  description: string
  instruction: string
  // globalInstruction: string
  tools: [...string]
  subagents: [...string]
}

agents: hack: {
  description: "A development agent to test out prompts, tools, and agents."
  instruction: "agents/hack.md"
  tools: [
    "cache_write",
    "cache_remove",
    "cache_file",
    "cache_dir",
  ]
}

agents: veggie: {
  description: string | *"Veggie, a general assistant helpful for any task"
  instruction: string | *"system/veggie.md"
  tools: [
    "cache_put",
    "cache_del",
  ]
}

agents: coding_context_provider: {
  description: string | *"Returns the relevant context from directory listings, file contents, and/or terminal history necessary to aid completing a task based on the query"
  instruction: string | *"agents/coding_context_provider.md"
  tools: [
    "cache_put",
    "cache_del",
    "fs_read",
    "fs_list",
    "fs_grep",
  ]
}

agents: coding_assistant: {
  description: string | *"A coding assistant for senior developers."
  instruction: string | *"agents/coding_assistant.md"
  tools: [
    // "@coding_context_provider",
    "cache_put",
    "cache_del",
    "fs_read",
    "fs_list",
    "fs_grep",
    "fs_write",
    "fs_edit",
    "fs_del",
  ]
}