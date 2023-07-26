package registry

#Versions: {
	hof: string
	cue: string
	go:  string

	// container-...
	runtime: string
	version: string
}

#Case: {
	// injected by harmony
	group: string
	name:  string
	vers:  #Versions
	// possibly pinned for specific reproducers?
	// maybe a separate directory for those as we find them.

	// configured by user, details about their project
	url:  string
	ref:  string
	type: "cli" | "pkg"
	scripts: [...string]
}

versions: #Versions

Registry: [g=string]: [n=string]: #Case & {group: g, name: n, vers: versions}
