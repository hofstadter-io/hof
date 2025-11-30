@extern(embed)

package embed

// todo, something like this under .veg will be default
// ... or better yet, we can pull a module in on init
// we also need to combine home & project agentic stuff

_i1: _ @embed(glob=*/*.md,type=text)
for path, content in _i1 { (path): content }

_i2: _ @embed(glob=*/*/*.md,type=text)
for path, content in _i2 { (path): content }

// _i3: _ @embed(glob=*/*/*/*.md,type=text)
// for path, content in _i3 { (path): content }
