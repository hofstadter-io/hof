DIRS=$(shell ls | grep -v Makefile)

dirs:
	echo $(DIRS)

BUILDS = $(addsuffix .build,$(DIRS))
.PHONY: build $(BUILDS)
build: $(BUILDS)
$(BUILDS):
	cd $(patsubst %.build,%,$@) && pnpm build
