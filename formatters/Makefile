include ../ci/make.inc

TOOLS := $(shell ls tools/)
.PHONY: tools
tools:
	@echo $(TOOLS)

# make image.*
IMAGES=$(addsuffix .image,$(TOOLS))
.PHONY: images
images: $(IMAGES)
$(IMAGES):
	docker build -t hofstadter/fmt-$(@:%.image=%):$(IMAGE_TAG) -f tools/$(@:%.image=%)/Dockerfile.debian tools/$(@:%.image=%)

MARCHS_DEB=$(addsuffix .march.debian,$(TOOLS))
marchs.debian: $(MARCHS_DEB)
$(MARCHS_DEB):
	docker buildx build -t hofstadter/fmt-$(@:%.march.debian=%):$(IMAGE_TAG)-debian \
		$(PUSH) --platform $(PLATFORMS) $(EXTRA) \
		tools/$(@:%.march.debian=%) -f tools/$(@:%.march.debian=%)/Dockerfile.debian
MARCHS_ALP=$(addsuffix .march.alpine,$(TOOLS))
marchs.alpine: $(MARCHS_ALP)
$(MARCHS_ALP):
	docker buildx build -t hofstadter/fmt-$(@:%.march.alpine=%):$(IMAGE_TAG)-alpine \
		$(PUSH) --platform $(PLATFORMS) $(EXTRA) \
		tools/$(@:%.march.alpine=%) -f tools/$(@:%.march.alpine=%)/Dockerfile.alpine

# make *.start
.PHONY: $(addsuffix .start,$(TOOLS))
$(addsuffix .start,$(TOOLS)):
	docker run -d -P --name hof-fmt-$(@:%.start=%) hofstadter/fmt-$(@:%.start=%):$(IMAGE_TAG)

# make *.stop
stop: $(addsuffix .stop,$(TOOLS))
$(addsuffix .stop,$(TOOLS)):
	docker rm -f hof-fmt-$(@:%.stop=%)

