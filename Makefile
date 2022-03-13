CUE_FILES  = $(shell find . -type f -name '*.cue' | grep -v 'cue.mod/pkg/' | sort)
GO_FILES  = $(shell find . -type f -name '*.go' | grep -v 'cue.mod/pkg/' | sort)
GHA_FILES  = $(shell ls .github/workflows/*.cue | sort)

# First command incase users forget to supply one
# cat self as help for simplicity
help:
	cat Makefile
.PHONY: help

fmt: cuefmt gofmt

cuefmt: $(CUE_FILES)
	@for f in $(CUE_FILES); do echo $$f; done

gofmt: $(GO_FILES)
	@for f in $(GO_FILES); do echo $$f; done

.PHONY: workflow
workflows = $(addprefix workflow_, $(GHA_FILES))
workflow: $(workflows)
$(workflows): workflow_%:
	@cue export --out yaml $(subst workflow_,,$@) > $(subst workflow_,,$(subst .cue,,$@)).yml
