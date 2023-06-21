include ../ci/make.inc

ifeq (, $(shell which docker))
	 $(error "Docker is not installed.")
endif

TOOLS := $(shell ls tools/)
.PHONY: tools
tools:
	@echo $(TOOLS)

.PHONY: buildx
buildx:
	-docker buildx rm mybuilder
	docker buildx create --name mybuilder --driver docker-container --bootstrap --use
	docker buildx inspect

# make image.*
IMAGES=$(addsuffix .image,$(TOOLS))
.PHONY: images
images: $(IMAGES)
$(IMAGES):
	docker build -t ghcr.io/hofstadter-io/fmt-$(@:%.image=%):$(FMTR_TAG) -f tools/$(@:%.image=%)/Dockerfile.debian tools/$(@:%.image=%)

MARCHS_DEB=$(addsuffix .march.debian,$(TOOLS))
marchs.debian: $(MARCHS_DEB)
$(MARCHS_DEB):
	docker buildx build -t ghcr.io/hofstadter-io/fmt-$(@:%.march.debian=%):$(FMTR_TAG) \
		$(PUSH) --platform $(PLATFORMS) $(EXTRA) \
		tools/$(@:%.march.debian=%) -f tools/$(@:%.march.debian=%)/Dockerfile.debian

# make *.start
.PHONY: $(addsuffix .start,$(TOOLS))
$(addsuffix .start,$(TOOLS)):
	docker run -d -P --name hof-fmt-$(@:%.start=%) ghcr.io/hofstadter-io/fmt-$(@:%.start=%):$(FMTR_TAG)

# make *.stop
stop: $(addsuffix .stop,$(TOOLS))
$(addsuffix .stop,$(TOOLS)):
	docker rm -f hof-fmt-$(@:%.stop=%)

