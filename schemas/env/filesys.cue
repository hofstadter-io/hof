package env

// like dagger.File
#File: Ref & {
	$kind:   "#file"
	path:    string
	source?: #Dir | #Container | #HostDir | #HostImage | #GitRepo
}

// this is creating a directory ref that we can do things with
#Dir: Ref & {
	$kind:   "#dir"
	source?: #Dir | #Container | #HostDir | #HostImage | #GitRepo
	path:    string
	include: [...string]
	exclude: [...string]
	gitignore: bool | *true
}

// like dagger.WithFile
File: Step & {
	$kind: "file"

	path: string

	// what if from another container? helper to make one of these from a Container?
	content?: string | #File | #HostFile // HMMM(A): should this just be file, or be container/image too?

	permissions?: int
	owner?:       string
	expand?:      bool
}

// this is including a directory in a container
Dir: Step & {
	$kind: "dir"
	// args
	path:    string
	source?: #Container | #Dir | #GitRepo | #HostDir | #HostImage // HMMM(B): or maybe this should just be dir kinds, make the user do an extra step? (nah, wouldn't have to with the SDK directly)
	// opts
	include?: [...string]
	exclude?: [...string]
	gitignore?: bool | *true
	owner?:     string
	expand?:    bool
}
