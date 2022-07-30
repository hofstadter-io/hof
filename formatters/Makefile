TAG := $(shell git tag --points-at HEAD | tr -d "\n")
TAG := $(if $(TAG),$(TAG),dirty)
tag:
	@echo $(TAG)

TOOLS := $(shell ls tools/)
tools:
	@echo $(TOOLS)

# make image.*
IMAGES=$(addsuffix .image,$(TOOLS))
images: $(IMAGES)
$(IMAGES):
	docker build -t hofstadter/fmt-$(@:%.image=%):$(TAG) tools/$(@:%.image=%)

# make *.run
$(addsuffix .run,$(TOOLS)):
	docker run -d -P --name hof-fmt-$(@:%.run=%) hofstadter/fmt-$(@:%.run=%):$(TAG)

# make *.stop
stop: $(addsuffix .stop,$(TOOLS))
$(addsuffix .stop,$(TOOLS)):
	docker rm -f hof-fmt-$(@:%.stop=%)

