# https://github.com/stevenvachon/broken-link-checker
BLC_EXCLUDES=--exclude 'https://github.com/hofstadter-io/hof/issues/new' \
	--exclude 'https://github.com/hofstadter-io/hof/edit'
blc.dev:
	blc -ro http://localhost:1313 ${BLC_EXCLUDES}
blc.prd:
	blc -ro https://docs.hofstadter.io ${BLC_EXCLUDES}


.PHONY: cuefmt cuefiles
cuefiles_all:
	find code/ -type f -name '*.cue' '!' -path '*/cue.mod/*' -print
cuefiles:
	find code/ -type f -name '*.cue' '!' -path '*/cue.mod/*' '!' -path '*/templates/*' '!' -path '*/partials/*' '!' -path '*/.hof/*' -print
cuefmt:
	find code/ -type f -name '*.cue' '!' -path '*/cue.mod/*' '!' -path '*/templates/*' '!' -path '*/partials/*' '!' -path '*/.hof/*' -exec cue fmt {} \;

.PHONY: gofmt gofiles
gofiles_all:
	find code/ -type f -name '*.go' '!' -path '*/cue.mod/*' -print
gofiles:
	find code/ -type f -name '*.go' '!' -path '*/cue.mod/*' '!' -path '*/templates/*' '!' -path '*/partials/*' '!' -path '*/.hof/*' -print
gofmt:
	find code/ -type f -name '*.go' '!' -path '*/cue.mod/*' '!' -path '*/templates/*' '!' -path '*/partials/*' '!' -path '*/.hof/*' -exec gofmt -w {} \;

