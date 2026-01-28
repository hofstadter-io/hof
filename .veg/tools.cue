package veg

import (
	"github.com/hofstadter-io/hof/.veg/embed" // puke...
)

tools: [n=string]: {
  @agentic(tool)
	name:        string | *n
	description: string
}

TOOLS=tools: {
	cache_put: description:    embed["tools/cache_put.md"]
	// cache_write: description:  embed["tools/cache_put.md"]
	// cache_edit: description:   embed["tools/cache_edit.md"]
	// cache_remove: description: embed["tools/cache_del.md"]
	cache_del: description:    embed["tools/cache_del.md"]

	fs_read: description:  embed["tools/fs_read.md"]
	fs_list: description:  embed["tools/fs_list.md"]
	fs_glob: description:  embed["tools/fs_glob.md"]
	fs_grep: description:  embed["tools/fs_grep.md"]

	fs_edit: description:  embed["tools/fs_edit.md"]
	fs_write: description: embed["tools/fs_write.md"]
	fs_del: description:   embed["tools/fs_del.md"]

	exec: description:   embed["tools/exec.md"]
}

toolsets: {
	[n=string]: { name: n }
	cache_only: {
		tools: [
			TOOLS.cache_put,
			TOOLS.cache_del,
		]
	}
	fs_query: {
		tools: [
			TOOLS.fs_read,
			TOOLS.fs_list,
			TOOLS.fs_glob,
			TOOLS.fs_grep,
		]
	}
	fs_mutate: {
		tools: [
			TOOLS.fs_edit,
			TOOLS.fs_write,
			TOOLS.fs_del,
		]
	}
}

mcp: {
	github: {
		uri: "https://api.githubcopilot.com/mcp/"
		envVar: "GITHUB_PAT"
	}
	tavily: {
		uri: "https://mcp.tavily.com/mcp/"
		envVar: "TAVILY_APIKEY"
	}
}
