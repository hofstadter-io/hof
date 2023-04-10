include ./hack/make/*.inc

COMMIT=$(shell git rev-parse --short HEAD | tr -d "\n")
TAG?=$(COMMIT)
PROJECT="hof-io--develop"

help:
	@cat Makefile

tag:
	@echo "tag: '$(TAG)'"

.PHONY: run lint test
run: dev
lint: fmt broken-link
test: verify

.PHONY: fmt broken-link highlight
fmt: cuefmt gofmt
broken-link: blc.dev
highlight: highlight-cue

# build the world and push to production
.PHONY: first prettier gen tools
first: deps prettier gen
gen: config.yaml extern examples highlight
prettier:
	hof fmt start prettier
