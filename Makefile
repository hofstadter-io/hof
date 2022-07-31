CUE_FILES  = $(shell find . -type f -name '*.cue' | grep -v 'cue.mod/pkg/' | sort)
GO_FILES  = $(shell find . -type f -name '*.go' | grep -v 'cue.mod/pkg/' | sort)
GHA_FILES  = $(shell ls .github/workflows/*.cue | sort)

# First command incase users forget to supply one
# cat self as help for simplicity
help:
	cat Makefile
.PHONY: help

# GitHub Actions workflows
.PHONY: workflow
workflows = $(addprefix workflow_, $(GHA_FILES))
workflow: $(workflows)
$(workflows): workflow_%:
	@cue export --out yaml $(subst workflow_,,$@) > $(subst workflow_,,$(subst .cue,,$@)).yml

fmt: cuefmt gofmt

.PHONY: cuefmt cuefiles
cuefiles:
	find . -type f -name '*.cue' '!' -path '*/cue.mod/*' '!' -path '*/templates/*' '!' -path '*/partials/*' '!' -path '*/.hof/*' -print
cuefmt:
	find . -type f -name '*.cue' '!' -path '*/cue.mod/*' '!' -path '*/templates/*' '!' -path '*/partials/*' '!' -path '*/.hof/*' -exec cue fmt {} \;

.PHONY: gofmt gofiles
gofiles:
	find . -type f -name '*.go' '!' -path '*/cue.mod/*' '!' -path '*/templates/*' '!' -path '*/partials/*' '!' -path '*/.hof/*' -print
gofmt:
	find . -type f -name '*.go' '!' -path '*/cue.mod/*' '!' -path '*/templates/*' '!' -path '*/partials/*' '!' -path '*/.hof/*' -exec gofmt -w {} \;

release:
	make -C formatters images
	cd cmd/hof && goreleaser --rm-dist -p 1
	make -C formatters push
