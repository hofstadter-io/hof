# run locally in dev mode
.PHONY: dev
dev: config.yaml
	@hugo serve --bind 0.0.0.0 --buildDrafts --buildFuture --disableFastRender

# run locally in prd mode
.PHONY: prd
stg: config.yaml
	@hugo serve --bind 0.0.0.0 --disableFastRender


