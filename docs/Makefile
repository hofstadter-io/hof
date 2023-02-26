include ./hack/make/*.inc

COMMIT  = $(shell git rev-parse --short HEAD | tr -d "\n")
PROJECT = "hof-io--develop"

help:
	@cat Makefile

.PHONY: run lint test
run: dev
lint: fmt broken-link
test: verify

.PHONY: fmt broken-link highlight
fmt: cuefmt gofmt
broken-link: blc.dev
highlight: highlight-cue

# build the world and push to production
.PHONY: first gen
first: deps gen
gen: config.yaml extern examples highlight

