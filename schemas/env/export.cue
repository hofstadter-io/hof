package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

#ExportDir: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "exportDir"
	}

	$kind: "#exportDir"
	name?:  string

	// where to place
	path: string

	// pieces that make up the bundled dir
	// sources:  [...#FileLike|#DirLike]
	sources: [...]

	// (1) filters
	include: [...string]
	exclude: [...string]
	gitignore: bool | *true

	// (2) path to select from the bundled dir
	trimPrefix?: string

	// (3) git-compatible patch to apply after bundling and selecting
	patch?:     string
	patchFile?: #FileLike

	// If true, then the host directory will be wiped clean before exporting so that it exactly matches the directory being exported; this means it will delete any files on the host that aren't in the exported dir. If false (the default), the contents of the directory will be merged with any existing contents of the host directory, leaving any existing files on the host that aren't in the exported directory alone.
	wipe: bool | *false
}

#ExportFile: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "exportFile"
	}

	$kind: "#exportFile"
	name?:  string
	path:  string
	file:  #File

	// If allowParentDirPath is true, the path argument can be a directory path, in which case the file will be created in that directory.
	allowParentDirPath?: bool
}

#ExportImageFile: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "exportImageFile"
	}

	$kind: "#exportImageFile"
	name?:  string
	path:  string
	tags: [...string]
	image: #Container
}

#ExportImage: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "exportImage"
	}

	$kind: "#exportImage"
	name?:  string
	reg?:  string
	tags: [...string]
	image: #Container
}

#PublishImage: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "publishImage"
	}

	$kind: "#publishImage"
	name?:  string
	reg:   string
	tags: [...string]
	image: #Container
}