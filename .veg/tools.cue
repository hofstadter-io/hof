package veg

import (
	"github.com/hofstadter-io/hof/.veg/embed" // puke...
)

tools: [n=string]: {
	name:        string | *n
	description: string
}

tools: {
	cache_put: description:    embed["tools/cache_put.md"]
	cache_write: description:  embed["tools/cache_put.md"]
	cache_edit: description:   embed["tools/cache_edit.md"]
	cache_remove: description: embed["tools/cache_del.md"]
	cache_del: description:    embed["tools/cache_del.md"]

	fs_read: description:  embed["tools/fs_read.md"]
	fs_list: description:  embed["tools/fs_list.md"]
	fs_grep: description:  embed["tools/fs_grep.md"]
	fs_edit: description:  embed["tools/fs_edit.md"]
	fs_write: description: embed["tools/fs_write.md"]
	fs_del: description:   embed["tools/fs_del.md"]

	exec: description:   embed["tools/exec.md"]
}
