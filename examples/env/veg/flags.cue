package veg

_flags: {
	repo: string | *"https://github.com/hofstadter-io/hof" @tag(repo)
	code: string | *"../../../"                            @tag(code)

	use: {
		lsp: bool | *false
	}

	ports: {
		gopls: int | *4000 @tag(ports_gopls)
		cupls: int | *4001 @tag(ports_cuepls)
	}
}
