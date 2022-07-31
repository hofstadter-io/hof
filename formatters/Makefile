DIFF := $(shell git diff)
TAG := $(shell git tag --points-at HEAD | tr -d "\n")

# overrides
TAG := $(if $(TAG),$(TAG),dirty)
TAG := $(if $(DIFF),dirty,$(TAG))

# informational
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

# make *.start
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

