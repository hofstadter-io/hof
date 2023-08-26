include ./ci/make.inc

CUE_FILES  = $(shell find . -type f -name '*.cue' | grep -v 'cue.mod/pkg/' | sort)
GO_FILES  = $(shell find . -type f -name '*.go' | grep -v 'cue.mod/pkg/' | sort)
GHA_FILES  = $(shell ls ci/gha/*.cue | sort)

# First command incase users forget to supply one
# cat self as help for simplicity
help:
	cat Makefile
.PHONY: help

.PHONY: github workflow devcontainer
github: workflow devcontainer
# GitHub Actions workflows
workflows = $(addprefix workflow_, $(GHA_FILES))
workflow: $(workflows)
$(workflows): workflow_%:
	@cue export --out yaml $(subst workflow_,,$@) -f -o $(subst ci/gha,.github/workflows,$(subst workflow_,,$(subst .cue,,$@))).yml
devcontainer: ci/devc/devcontainer.cue
	@cue export ci/devc/devcontainer.cue -f -o .devcontainer/devcontainer.json

.PHONY: hack 
hack:
	CGO_ENABLED=0 go install ./hack

.PHONY: hof
hof:
	CGO_ENABLED=0 go install ./cmd/hof

.PHONY: race
race:
	go install -race ./cmd/hof

.PHONY: hof.build
hof.build:
	CGO_ENABLED=0 go build -o hof ./cmd/hof

.PHONY: deps.check deps.update
deps.check:
	go list -u -m all | grep '\[.*\]'
deps.update:
	go get -u ./...

# formatter images
.PHONY: formatters fmtr.release
formatters:
	make -C formatters images
fmtr.release:
	make -C formatters marchs

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

goreleaser.yml: cmd/hof/goreleaser.cue
	cue export cmd/hof/goreleaser.cue -f -o cmd/hof/goreleaser.yml
snapshot: goreleaser.yml
	cd cmd/hof && goreleaser release -f goreleaser.yml --rm-dist -p 1 --snapshot
release: goreleaser.yml
	make -C formatters marchs
	cd cmd/hof && goreleaser release -f goreleaser.yml --rm-dist -p 1

workdir.clean:
	find . -type d -name '.workdir' -exec rm -rf {} \;
