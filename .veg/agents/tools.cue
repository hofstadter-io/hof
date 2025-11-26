package agents

tools: [n=string]: {
	name:        string | *n
	description: string
}

tools: {
	cache_put: description:    instructions.tools.cache_write
	cache_write: description:  instructions.tools.cache_write
	cache_edit: description:   instructions.tools.cache_edit
	cache_remove: description: instructions.tools.cache_remove
	cache_del: description:    instructions.tools.cache_remove

	fs_read: description:  instructions.tools.fs_read
	fs_list: description:  instructions.tools.fs_list
	fs_grep: description:  instructions.tools.fs_grep
	fs_edit: description:  instructions.tools.fs_edit
	fs_write: description: instructions.tools.fs_write
	fs_del: description:   instructions.tools.fs_del
}
