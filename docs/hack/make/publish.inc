build: hugo docker

config.yaml: config.cue
	cue export config.cue --out yaml --outfile config.yaml --force

.PHONY: hugo hugo.prod
hugo: hugo.prod
hugo.prod: config.yaml
	@rm -rf dist
	@hugo --baseURL https://docs.hofstadter.io/ -d dist
hugo.next: config.yaml
	@rm -rf dist
	@hugo --baseURL https://next.hofstadter.io/ -d dist

.PHONY: docker
docker: image

.PHONY: image image.next
image:
	@docker build --no-cache -f ci/Dockerfile -t us.gcr.io/$(PROJECT)/docs.hofstadter.io:$(TAG) .
image.next:
	@docker build --no-cache -f ci/Dockerfile -t us.gcr.io/$(PROJECT)/next.hofstadter.io:$(TAG) .

.PHONY: nginx
nginx:
	@docker run --rm -it -p 8080:80 -e GA_MP_APIKEY=${GA_MP_APIKEY} --name hof-docs us.gcr.io/$(PROJECT)/docs.hofstadter.io:$(TAG)

.PHONY: push
push:
	@docker push us.gcr.io/$(PROJECT)/docs.hofstadter.io:$(TAG)
push.next:
	@docker push us.gcr.io/$(PROJECT)/next.hofstadter.io:$(TAG)

.PHONY: deploy.prod deploy.prod.view deploy.prod.vet deploy.prod.all
deploy.prod.all: build push deploy.prod
deploy.prod:
	cue export ci/cuelm.cue -e Install \
		-t ga_mp_apikey=$(GA_MP_APIKEY) \
		-t version=$(TAG) \
		| kubectl apply -f -
deploy.prod.view:
	cue export ci/cuelm.cue -e Install \
		-t ga_mp_apikey=$(GA_MP_APIKEY) \
		-t version=$(TAG)
deploy.prod.vet:
	cue export ci/cuelm.cue -e Install \
		-t ga_mp_apikey=$(GA_MP_APIKEY) \
		-t version=$(TAG) \
		| kubectl apply -f - --dry-run=server

.PHONY: deploy.next deploy.next.view deploy.next.vet deploy.next.all
deploy.next.all: hugo.next image.next push.next deploy.next
deploy.next:
	cue export ci/cuelm.cue -e Install \
		-t ga_mp_apikey=$(GA_MP_APIKEY) \
		-t version=$(TAG) \
		-t name=hof-next-docs \
		-t domain="next.hofstadter.io" \
		| kubectl apply -f -
deploy.next.view:
	cue export ci/cuelm.cue -e Install \
		-t ga_mp_apikey=$(GA_MP_APIKEY) \
		-t version=$(TAG) \
		-t name=hof-next-docs\
		-t domain="next.hofstadter.io"
deploy.next.vet:
	cue export ci/cuelm.cue -e Install \
		-t ga_mp_apikey=$(GA_MP_APIKEY) \
		-t version=$(TAG) \
		-t name=hof-next-docs\
		-t domain="next.hofstadter.io" \
		| kubectl apply -f - --dry-run=server
