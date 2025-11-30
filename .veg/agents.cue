package veg

import (
	"github.com/hofstadter-io/hof/.veg/embed" // puke...
)

// embed the whole package (flat map of path->content)
embeds: embed
embedDir: "./.veg/embed"

agents: [n=string]: {
  name: string | *n
  model: string | *"gemini-2.5-flash"
  description: string
  instruction: string
  // globalInstruction: string
  tools: [...string]
  subagents: [...string]
}

agents: veggie: {
  description: string | *"Veggie, a general assistant helpful for any task"
  instruction: string | *"system/veggie.md"
  tools: [
    "cache_put",
    "cache_del",
  ]
}

agents: coding_assistant: {
  description: string | *"Veggie Code, a sophisticated assistant for senior developers."
  instruction: string | *"agents/coding_assistant.md"
  runenv: "golang"
  tools: [
    "cache_put",
    "cache_del",
    "fs_read",
    "fs_list",
    "fs_grep",
    "fs_write",
    "fs_edit",
    "fs_del",
    "exec",
    // "@coding_context_provider",
  ]
}

agents: coding_assro: {
  description: string | *"Veggie Code, a sophisticated assistant for senior developers."
  instruction: string | *"agents/coding_assro.md"
  runenv: "golang"
  tools: [
    "cache_put",
    "cache_del",
    "fs_read",
    "fs_list",
    "fs_grep",
    "exec",
    // "@coding_context_provider",
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

agents: deepc: {
  description: "A deep research agent specializing in code base analysis to iteratively search, evaluate, summarize, and build well cited research."
  instruction: "agents/deepr.md"
  tools: [
    "cache_put",
    "cache_del",
    "fs_read",
    "fs_list",
    "fs_grep",
    "fs_write",
    "search",
    "fetch",
  ]
}

agents: deepr: {
  description: "A deep research agent to iteratively search, evaluate, summarize, and build well cited reports."
  instruction: "agents/deepr.md"
  tools: [
    "cache_put",
    "cache_del",
    "fs_read",
    "fs_list",
    "fs_grep",
    "fs_write",
    "search",
    "fetch",
  ]
}

agents: fetch: {
  description: "A web crawling agent."
  instruction: "agents/deepr.md"
  tools: [
    "cache_put",
    "cache_del",
    "search",
    "fetch",
  ]
}

agents: hack: {
  description: "A development agent to test out prompts, tools, and agents."
  instruction: "agents/hack.md"
  tools: [
    "cache_put",
    "cache_del",
    "fs_read",
    "fs_list",
    "fs_grep",
  ]
}
