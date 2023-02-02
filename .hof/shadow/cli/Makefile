CUE_FILES  = $(shell find . -type f -name '*.cue' | grep -v 'cue.mod/pkg/' | sort)
GO_FILES  = $(shell find . -type f -name '*.go' | grep -v 'cue.mod/pkg/' | sort)
# GHA_FILES  = $(shell ls .github/workflows/*.cue | sort)

# First command incase users forget to supply one
# cat self as help for simplicity
help:
	cat Makefile
.PHONY: help

# GitHub Actions workflows
workflows = $(addprefix workflow_, $(GHA_FILES))
workflow: $(workflows)
$(workflows): workflow_%:
	@cue export --out yaml $(subst workflow_,,$@) > $(subst workflow_,,$(subst .cue,,$@)).yml
.PHONY: workflow

deps:
	hof mod vendor cue
	go mod tidy
.PHONY: deps

build:
	CGO_ENABLED=0 go build -o hof ./cmd/hof
.PHONY: build

install:
	CGO_ENABLED=0 go install ./cmd/hof
.PHONY: install

fmt: fmt.go fmt.cue
fmt.go:
	find . -type f -name '*.go' '!' -path '*/cue.mod/*' '!' -path '*/templates/*' '!' -path '*/partials/*' '!' -path '*/.hof/*' -exec gofmt -w {} \;
.PHONY: fmt.go
fmt.cue:
	find . -type f -name '*.cue' '!' -path '*/cue.mod/*' '!' -path '*/templates/*' '!' -path '*/partials/*' '!' -path '*/.hof/*' -exec cue fmt {} \;
.PHONY: fmt.cue

release:
	cd cmd/hof && goreleaser --rm-dist -p 1
.PHONY: release

snapshot:
	cd cmd/hof && goreleaser --rm-dist -p 1 --snapshot
.PHONY: snapshot

