.PHONY: dev stg codespace

dev: config.yaml
	@hugo serve --bind 0.0.0.0 --buildDrafts --buildFuture --disableFastRender

stg: config.yaml
	@hugo serve --bind 0.0.0.0 --disableFastRender


codespace: config.yaml
	@hugo serve --baseUrl=/ --appendPort=false --bind 0.0.0.0 --buildDrafts --buildFuture --disableFastRender

