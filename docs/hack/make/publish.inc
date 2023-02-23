build: hugo docker
release: push deploy

config.yaml: config.cue
	cue export config.cue --out yaml --outfile config.yaml --force

.PHONY: hugo
hugo: config.yaml
	@rm -rf dist
	@hugo --baseURL https://docs.hofstadter.io/ -d dist

.PHONY: docker
docker: image push

.PHONY: image
image:
	@docker build --no-cache -f ci/Dockerfile -t us.gcr.io/$(PROJECT)/docs.hofstadter.io:$(TAG) .

.PHONY: nginx
nginx:
	@docker run --rm -it -p 8080:80 --name hof-docs us.gcr.io/$(PROJECT)/docs.hofstadter.io:$(TAG)

.PHONY: push
push:
	@docker push us.gcr.io/$(PROJECT)/docs.hofstadter.io:$(TAG)

.PHONY: deploy
deploy:
	@cue export ci/cuelm.cue -t version=$(TAG) -e Install | kubectl apply -f -

.PHONY: deploy-view
deploy-view:
	@cue export ci/cuelm.cue -t version=$(TAG) -e Install




