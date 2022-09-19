# git info
GIT_ROOT   :=$(shell git rev-parse --show-toplevel)
GIT_COMMIT :=$(shell git log -1 --pretty=format:"%h")
GIT_BRANCH :=$(shell git rev-parse --abbrev-ref HEAD | tr -d "\n")
GIT_TAG    :=$(shell git tag --points-at HEAD | tr -d "\n")
GIT_DIFF   :=$(shell git diff)
GIT_DIRTY  :=$(if $(GIT_DIFF),-dirty,)

# image tag setup
IMAGE_TAG :=$(GIT_COMMIT)
IMAGE_TAG :=$(if $(GIT_TAG),$(GIT_TAG),$(IMAGE_TAG))
IMAGE_TAG :=$(IMAGE_TAG)$(GIT_DIRTY)
FMTR_TAG  :=$(if $(GIT_TAG),$(GIT_TAG),dirty)
FMTR_TAG  :=$(if $(GIT_DIFF),dirty,$(FMTR_TAG))

# informational
info:
	@echo "GIT_ROOT:    $(GIT_ROOT)"
	@echo "GIT_COMMIT:  $(GIT_COMMIT)"
	@echo "GIT_BRANCH:  $(GIT_BRANCH)"
	@echo "GIT_TAG:     $(GIT_TAG)"
	@echo "GIT_DIRTY:   $(GIT_DIRTY)"
	@echo "IMAGE_TAG:   $(IMAGE_TAG)"
	@echo "FMTR_TAG:    $(FMTR_TAG)"

PUSH=--push
PLATFORMS=linux/amd64,linux/arm64
EXTRA=--progress plain --no-cache

.PHONY: marchs
marchs: marchs.debian marchs.alpine

# buildx / buildkit related
buildx.setup:
	docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
buildx.check:
	docker buildx ls
buildx.start:
	docker buildx create --name builder --driver docker-container --use
	docker buildx inspect --bootstrap
buildx.stop:
	docker buildx rm builder

