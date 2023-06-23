include ../ci/make.inc

ifneq (, $(shell which docker))
	export RUNTIME:=docker
else ifneq (, $(shell which podman))
	export RUNTIME:=podman
else
	 $(error "No OCI runtime detected. Please install Docker or Podman.")
endif

TOOLS := $(shell ls tools/)
.PHONY: tools
tools:
	@echo $(TOOLS)

.PHONY: buildx
buildx:
	${RUNTIME} buildx rm mybuilder
	${RUNTIME} buildx create --name mybuilder --driver docker-container --bootstrap --use
	${RUNTIME} buildx inspect

# make image.*
IMAGES=$(addsuffix .image,$(TOOLS))
.PHONY: images
images: $(IMAGES)
$(IMAGES):
	${RUNTIME} build -t ghcr.io/hofstadter-io/fmt-$(@:%.image=%):$(FMTR_TAG) -f tools/$(@:%.image=%)/Dockerfile.debian tools/$(@:%.image=%)

MARCHS_DEB=$(addsuffix .march.debian,$(TOOLS))
marchs.debian: $(MARCHS_DEB)
$(MARCHS_DEB):
	${RUNTIME} buildx build -t ghcr.io/hofstadter-io/fmt-$(@:%.march.debian=%):$(FMTR_TAG) \
		$(PUSH) --platform $(PLATFORMS) $(EXTRA) \
		tools/$(@:%.march.debian=%) -f tools/$(@:%.march.debian=%)/Dockerfile.debian

# make *.start
.PHONY: $(addsuffix .start,$(TOOLS))
$(addsuffix .start,$(TOOLS)):
	${RUNTIME} run -d -P --name hof-fmt-$(@:%.start=%) ghcr.io/hofstadter-io/fmt-$(@:%.start=%):$(FMTR_TAG)

# make *.stop
stop: $(addsuffix .stop,$(TOOLS))
$(addsuffix .stop,$(TOOLS)):
	${RUNTIME} rm -f hof-fmt-$(@:%.stop=%)

