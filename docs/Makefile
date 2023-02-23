include ./hack/make/*.inc

TAG        = $(shell git rev-parse --short HEAD | tr -d "\n")
PROJECT    = "hof-io--develop"

help:
	@cat Makefile

.PHONY: run lint test
run: dev
lint: fmt broken-link
test: verify

.PHONY: fmt broken-link code highlight
fmt: cuefmt gofmt
broken-link: blc.dev
examples: gen
highlight: highlight-cue

# build the world and push to production
.PHONY: first setup publish
first: deps setup
setup: config.yaml extern code highlight
publish: setup build deploy

