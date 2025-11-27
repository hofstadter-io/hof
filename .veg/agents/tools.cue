package agents

tools: [n=string]: {
	name:        string | *n
	description: string
}

tools: {
	cache_put: description:    instructions["tools/cache_write.md"]
	cache_write: description:  instructions["tools/cache_write.md"]
	cache_edit: description:   instructions["tools/cache_edit.md"]
	cache_remove: description: instructions["tools/cache_remove.md"]
	cache_del: description:    instructions["tools/cache_remove.md"]

	fs_read: description:  instructions["tools/fs_read.md"]
	fs_list: description:  instructions["tools/fs_list.md"]
	fs_grep: description:  instructions["tools/fs_grep.md"]
	fs_edit: description:  instructions["tools/fs_edit.md"]
	fs_write: description: instructions["tools/fs_write.md"]
	fs_del: description:   instructions["tools/fs_write.md"]
}
