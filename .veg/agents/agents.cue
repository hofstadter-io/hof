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

agents: general_assistant: {
  description: string | *"A general assistant helpful for any task"
  instruction: string | *instructions.agents.general_assistant
  tools: []
}

agents: file_system_context_provider: {
  description: string | *"Returns file contents and/or directory listings based on the query"
  instruction: string | *instructions.agents.file_system_context_provider
  tools: [
    "directory_tree",
    "list_directory",
    "grep_regexp",
    "read_file",
  ]
}

agents: coding_context_provider: {
  description: string | *"Returns the relevant context from directory listings, file contents, and/or terminal history necessary to aid completing a task based on the query"
  instruction: string | *instructions.agents.coding_context_provider
  tools: [
    "directory_tree",
    "list_directory",
    "grep_regexp",
    "read_file",
  ]
}

agents: coding_assistant: {
  description: string | *"A coding assistant for senior developers."
  instruction: string | *instructions.agents.coding_assistant
  tools: [
    "@coding_context_provider",
    "directory_tree",
    "list_directory",
    "grep_regexp",
    "read_file",
  ]
}