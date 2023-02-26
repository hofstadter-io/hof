.PHONY: gen
gen:
	cd simple-server && hof mod link && hof gen ./examples/
	cd data-layer    && hof mod link && hof gen ./examples/
	cd web-client    && hof mod link && hof gen ./examples/

test: gen
	git diff --exit-code
