include ../ci/make.inc

ifneq (, $(shell which docker))
	export CONTAINER_TOOL=docker
else ifneq (, $(shell which podman))
	export CONTAINER_TOOL=podman
else
	 $(error "No container tool found. Install Docker or Podman to continue.")
endif

TOOLS := $(shell ls tools/)
.PHONY: tools
tools:
	@echo $(TOOLS)

.PHONY: buildx
buildx:
	-${CONTAINER_TOOL} buildx rm mybuilder
	${CONTAINER_TOOL} buildx create --name mybuilder --driver docker-container --bootstrap --use
	${CONTAINER_TOOL} buildx inspect

# make image.*
IMAGES=$(addsuffix .image,$(TOOLS))
.PHONY: images
images: $(IMAGES)
$(IMAGES):
	${CONTAINER_TOOL} build -t ghcr.io/hofstadter-io/fmt-$(@:%.image=%):$(FMTR_TAG) -f tools/$(@:%.image=%)/Dockerfile.debian tools/$(@:%.image=%)

MARCHS_DEB=$(addsuffix .march.debian,$(TOOLS))
marchs.debian: $(MARCHS_DEB)
$(MARCHS_DEB):
	${CONTAINER_TOOL} buildx build -t ghcr.io/hofstadter-io/fmt-$(@:%.march.debian=%):$(FMTR_TAG) \
		$(PUSH) --platform $(PLATFORMS) $(EXTRA) \
		tools/$(@:%.march.debian=%) -f tools/$(@:%.march.debian=%)/Dockerfile.debian

# make *.start
.PHONY: $(addsuffix .start,$(TOOLS))
$(addsuffix .start,$(TOOLS)):
	${CONTAINER_TOOL} run -d -P --name hof-fmt-$(@:%.start=%) ghcr.io/hofstadter-io/fmt-$(@:%.start=%):$(FMTR_TAG)

# make *.stop
stop: $(addsuffix .stop,$(TOOLS))
$(addsuffix .stop,$(TOOLS)):
	${CONTAINER_TOOL} rm -f hof-fmt-$(@:%.stop=%)

