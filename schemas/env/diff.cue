package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

// #Changes calcs the changeset between two directories
// is a: *dagger.Changeset
#Changes: Ref & {
	schemas.Hof
	#hof: env: {
		root: true 
		kind: "changes"
	}
	$kind: "#changes"

  // #DirLike | #ImageLike
  prev: _
  next: _
}

// Applies a changeset to a container
Changes: Step & {
	$kind: "changes"

  change: #Changes
}

// create a #File from #Changes or git-like patch
#PatchFile: Ref & {
	schemas.Hof
	#hof: env: {
		root: true 
		kind: "patchFile"
	}
	$kind: "#patchFile"
  
	name?:  string
  source: string | #Changes
}

// patch a #Container with a git-like patch or #Changes
Patch: Step & {
	$kind: "patch"
  source: string | #Changes
	basepath?: string
}

// patch a #Container with a #PatchFile
PatchFile: Step & {
	$kind: "patchFile"
  source: #PatchFile
	basepath?: string
}

// #HostPatch

// #Shouldi resolves then or else
//   based on a diff and patterns
#Shouldi: Ref & {
	schemas.Hof
	#hof: env: {
		root: true 
		kind: "shouldi"
	}
	$kind: "#shouldi"
  
  // changes to match against
  changes!: #Changes

  // patterns for matching
  include: [...string]
  exclude: [...string]

  // just do it!
  force?: bool

  // what to do
  then!: _
  else?: _
}