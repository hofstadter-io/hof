package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

// a file ref that can be used within CUE
// is a: *dagger.File
#File: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "file"
	}

	$kind: "#file"
	name:  string | *path
	path!: string
	trimPrefix: string | *""

	// actual, import env/rrr:env to enforce, performance penalty included
	// source: #Dir | #Container | #HostDir | #HostImage | #GitRepo
	source!: _
}

// step that adds a file to a container
// is a: dagger.WithFile
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

// a dir ref that can be used within CUE
// is a: *dagger.Directory
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
	trimPrefix: string | *""

	// (3) git-compatible patch to apply after bundling and selecting
	// perhaps this gets moved out, or updated and kept for convenience
	patch?:     string | #Changes
	patchFile?: #PatchFile
}

// step that adds a dir to a container
// is a: dagger.WithDirectory
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
	trimPrefix?: string
	gitignore?: bool | *true
	owner?:     string
	expand?:    bool

	// maybe patch stuff here too? as a convenience
	patch?:     string | #Changes
	patchFile?: #PatchFile
}


// use the RootFS of a container as a dir ref that can be used within CUE
// is a: *dagger.Directory
#RootFS: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "rootfs"
	}

	$kind: "#rootfs"

	// source: #ImageLike
	source: _
}

// step that sets the RootFS of a container to the source dir
RootFS: Step & {
	$kind: "rootfs"

	// source: #DirLike
	source: _
}