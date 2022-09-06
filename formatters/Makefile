DIFF := $(shell git diff)
TAG := $(shell git tag --points-at HEAD | tr -d "\n")

# overrides
TAG := $(if $(TAG),$(TAG),dirty)
TAG := $(if $(DIFF),dirty,$(TAG))

# informational
tag:
	@echo $(TAG)

TOOLS := $(shell ls tools/)
.PHONY: tools
tools:
	@echo $(TOOLS)

# make image.*
IMAGES=$(addsuffix .image,$(TOOLS))
.PHONY: images
images: $(IMAGES)
$(IMAGES):
	docker build -t hofstadter/fmt-$(@:%.image=%):$(TAG) -f tools/$(@:%.image=%)/Dockerfile.debian tools/$(@:%.image=%)


# PLATFORMS=linux/amd64
PLATFORMS=linux/arm64
# PLATFORMS=linux/amd64,linux/arm64
# PLATFORMS=linux/amd64,linux/arm64,linux/arm/v7

PUSH=--push
PUSH=

EXTRA=--progress plain --no-cache

.PHONY: marchs
marchs: marchs.debian marchs.alpine

MARCHS_DEB=$(addsuffix .march.debian,$(TOOLS))
marchs.debian: $(MARCHS_DEB)
$(MARCHS_DEB):
	docker buildx build -t hofstadter/fmt-$(@:%.march.debian=%):$(TAG)-debian \
		$(PUSH) --platform $(PLATFORMS) $(EXTRA) \
		tools/$(@:%.march.debian=%) -f tools/$(@:%.march.debian=%)/Dockerfile.debian
MARCHS_ALP=$(addsuffix .march.alpine,$(TOOLS))
marchs.alpine: $(MARCHS_ALP)
$(MARCHS_ALP):
	docker buildx build -t hofstadter/fmt-$(@:%.march.alpine=%):$(TAG)-alpine \
		$(PUSH) --platform $(PLATFORMS) $(EXTRA) \
		tools/$(@:%.march.alpine=%) -f tools/$(@:%.march.alpine=%)/Dockerfile.alpine

# make *.start
.PHONY: $(addsuffix .start,$(TOOLS))
$(addsuffix .start,$(TOOLS)):
	docker run -d -P --name hof-fmt-$(@:%.start=%) hofstadter/fmt-$(@:%.start=%):$(TAG)

# make *.stop
stop: $(addsuffix .stop,$(TOOLS))
$(addsuffix .stop,$(TOOLS)):
	docker rm -f hof-fmt-$(@:%.stop=%)

# make *.push
push: $(addsuffix .push,$(TOOLS))
$(addsuffix .push,$(TOOLS)):
	docker push hofstadter/fmt-$(@:%.push=%):$(TAG)

