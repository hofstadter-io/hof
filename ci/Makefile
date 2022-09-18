include ./make.inc

PUSH=
PLATFORMS=linux/amd64
EXTRA=--progress plain

marchs.debian:
	docker buildx build -t hofstadter/hof:$(IMAGE_TAG)-debian \
		$(PUSH) --platform $(PLATFORMS) $(EXTRA) \
		$(GIT_ROOT) -f $(GIT_ROOT)/ci/hof/docker/Dockerfile.debian
marchs.alpine:
	docker buildx build -t hofstadter/hof:$(IMAGE_TAG)-alpine \
		$(PUSH) --platform $(PLATFORMS) $(EXTRA) \
		$(GIT_ROOT) -f $(GIT_ROOT)/ci/hof/docker/Dockerfile.alpine
