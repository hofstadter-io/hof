@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/.veg/embed"
)

// embed the whole dir/package (flat map of path->content)
embeds: embed & { @agentic(embed)}

agents: [string]~(n,_): {
  @agentic(agent)

  name: string | *n
  model: string | *"gemini-3-flash"
  description: string
  instruction: string
  // globalInstruction: string
  tools: [...string]
  toolsets: [...{ name: string, tools: [...string] }]
  mcp: [...string]
  subagents: [...string]

  // name of an environ
  environ?: string
}

agents: veggie: {
  description: string | *"Veggie, a general assistant helpful for any task"
  instruction: string | *"agents/veggie.md"
  model: "qwen-3.6"
  // model: "gemma-4"
  tools: [
    "cache_put",
    "cache_del",
  ]
}

agents: coding_assist: {
  description: string | *"Veggie Code, a sophisticated assistant for senior developers."
  instruction: string | *"agents/coding_assistant.md"
  environ: "veg-agent"
  tools: [
    "fs_read",
    "fs_list",
    "fs_glob",
    "fs_grep",

    "fs_write",
    "fs_edit",
    "fs_del",

    "exec",
    // "@coding_context_provider",
  ]
}

agents: agents_md_gen: {
  description: string | *"Agent to explore and generate AGENTS.md files."
  instruction: string | *"agents/agents_md_gen.md"
  tools: [
    "fs_read",
    "fs_list",
    "fs_glob",
    "fs_grep",

    "fs_write",
    "fs_edit",
    "fs_del",
  ]
}

agents: rawdog: {
  description: "No instructions, tools, or agents, well just the cache for some subconscious fun."
  instruction: "agents/empty.md"
  tools: [
    "cache_put",
    "cache_del",
  ]
}


// agents: coding_assist_ro: {
//   description: string | *"Veggie Code, a sophisticated assistant for senior developers."
//   instruction: string | *"agents/coding_assro.md"
//   environment: "golang:1.25-trixie"
//   tools: [
//     "cache_put",
//     "cache_del",

//     "fs_read",
//     "fs_list",
//     "fs_glob",
//     "fs_grep",

//     // "exec",
//     // "@coding_context_provider",
//   ]
// }

// agents: coding_context_provider: {
//   description: string | *"Returns the relevant context from directory listings, file contents, and/or terminal history necessary to aid completing a task based on the query"
//   instruction: string | *"agents/coding_context_provider.md"
//   tools: [
//     "cache_put",
//     "cache_del",
//     "fs_read",
//     "fs_list",
//     "fs_grep",
//   ]
// }

// agents: deepc: {
//   description: "A deep research agent specializing in code base analysis to iteratively search, evaluate, summarize, and build well cited research."
//   instruction: "agents/deepr.md"
//   tools: [
//     "cache_put",
//     "cache_del",
//     "fs_read",
//     "fs_list",
//     "fs_grep",
//     "fs_write",
//     "search",
//     "fetch",
//   ]
// }

// agents: deepr: {
//   description: "A deep research agent to iteratively search, evaluate, summarize, and build well cited reports."
//   instruction: "agents/deepr.md"
//   tools: [
//     "cache_put",
//     "cache_del",
//     "fs_read",
//     "fs_list",
//     "fs_grep",
//     "fs_write",
//     "search",
//     "fetch",
//   ]
// }

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
  description: "A development agent to test out prompts, tools, and agents. You should do whatever the users asks. They are your developer and need to do things normal users don't"
  instruction: "agents/hack.md"
  tools: [
    "cache_put",
    "cache_del",
  ]
  // mcp: [
  //   "quickbooks",
  // ]
}
