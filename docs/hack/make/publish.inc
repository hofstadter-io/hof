config.yaml: config.cue
	cue export config.cue --out yaml --outfile config.yaml --force

.PHONT: nginx.next nginx.prod
nginx.next:
	@docker run --rm -it -p 8080:80 -e GA_MP_APIKEY=${GA_MP_APIKEY} --name hof-docs us.gcr.io/$(PROJECT)/next.hofstadter.io:$(TAG)
nginx.prod:
	@docker run --rm -it -p 8080:80 -e GA_MP_APIKEY=${GA_MP_APIKEY} --name hof-docs us.gcr.io/$(PROJECT)/docs.hofstadter.io:$(TAG)


.PHONY: all.prod hugo.prod image.prod nginx.prod deploy.prod deploy.prod.view deploy.prod.vet
all.prod: hugo.prod image.prod push.prod deploy.prod

hugo.prod: config.yaml
	@rm -rf dist
	@hugo --baseURL https://docs.hofstadter.io/ -d dist
image.prod:
	@docker build --no-cache -f ci/Dockerfile -t us.gcr.io/$(PROJECT)/docs.hofstadter.io:$(TAG) .
push.prod:
	@docker push us.gcr.io/$(PROJECT)/docs.hofstadter.io:$(TAG)

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


.PHONY: all.next hugo.next image.next nginx.next deploy.next deploy.next.view deploy.next.vet
all.next: hugo.next image.next push.next deploy.next

hugo.next: config.yaml
	@rm -rf dist
	@hugo --baseURL https://next.hofstadter.io/ -d dist
image.next:
	@docker build --no-cache -f ci/Dockerfile -t us.gcr.io/$(PROJECT)/next.hofstadter.io:$(TAG) .
push.next:
	@docker push us.gcr.io/$(PROJECT)/next.hofstadter.io:$(TAG)

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
