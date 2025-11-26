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
  instruction: instructions.agents.hack
  tools: [
    "cache_write",
    "cache_remove",
    "cache_file",
    "cache_dir",
  ]
}

agents: veggie: {
  description: string | *"Veggie, a general assistant helpful for any task"
  instruction: string | *instructions.system.veggie
  tools: []
}

agents: coding_context_provider: {
  description: string | *"Returns the relevant context from directory listings, file contents, and/or terminal history necessary to aid completing a task based on the query"
  instruction: string | *instructions.agents.coding_context_provider
  tools: [
    "cache_write",
    "cache_remove",
    "cache_grep",
    "cache_file",
    "cache_dir",
  ]
}

agents: coding_assistant: {
  description: string | *"A coding assistant for senior developers."
  instruction: string | *instructions.agents.coding_assistant
  tools: [
    // "@coding_context_provider",
    "cache_write",
    "cache_edit",
    "cache_remove",
    "cache_grep",
    "cache_file",
    "cache_dir",
  ]
}