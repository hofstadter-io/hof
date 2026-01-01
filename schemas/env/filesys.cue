package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

// like dagger.File
#File: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "file"
	}

	$kind: "#file"
	name:  string | *path
	path!: string

	// actual, import env/rrr:env to enforce, performance penalty included
	// source: #Dir | #Container | #HostDir | #HostImage | #GitRepo
	source!: _
}

// this is creating a directory ref that we can do things with
#Dir: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "dir"
	}

	$kind: "#dir"
	name:  string | *path

	// where to place
	path: string | *"."

	// pieces that make up the bundled dir
	// sources:  [...#FileLike|#DirLike]
	sources: [...]

	// (1) filters
	include: [...string]
	exclude: [...string]
	gitignore: bool | *true

	// (2) path to select from the bundled dir
	bundlePath: string | *"/"

	// (3) git-compatible patch to apply after bundling and selecting
	patch?:     string
	patchFile?: #FileLike
}

// like dagger.WithFile
File: Step & {
	$kind: "file"

	path!: string

	// actual, import env/rrr:env to enforce, performance penalty included
	// content: string | #File | #HostFile // HMMM(A): should this just be file, or be container/image too?
	content!: _

	permissions?: int
	owner?:       string
	expand?:      bool
}

// this is including a directory in a container
Dir: Step & {
	$kind: "dir"
	// args
	path: string | *"."

	// actual, import env/rrr:env to enforce, performance penalty included
	// source: #Container | #Dir | #GitRepo | #HostDir | #HostImage // HMMM(B): or maybe this should just be dir kinds, make the user do an extra step? (nah, wouldn't have to with the SDK directly)
	source!: _

	// opts
	include?: [...string]
	exclude?: [...string]
	gitignore?: bool | *true
	owner?:     string
	expand?:    bool
}
