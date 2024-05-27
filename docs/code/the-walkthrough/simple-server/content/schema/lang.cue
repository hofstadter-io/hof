Server: {
	// ...

	// The project's git repo
	GitRepo: string

	// We need to know the module for Go
	//   we can default to another field
	GoModule: string | *GitRepo
}
