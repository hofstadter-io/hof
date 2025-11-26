package agents

tools: [n=string]: {
  name: string | *n
  description: string
}

tools: {
	cache_write: description:  instructions.tools.cache_write
	cache_edit: description:   instructions.tools.cache_edit
	cache_remove: description: instructions.tools.cache_remove
	cache_grep: description:   instructions.tools.cache_grep
	cache_file: description:   instructions.tools.cache_file
	cache_dir: description:    instructions.tools.cache_dir
}
