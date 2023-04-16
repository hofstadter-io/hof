.PHONY: watch.adhoc watch.gen
watch.adhoc:
	hof gen data.cue -T min.txt --watch

watch.gen:
	pushd gen
	hof gen --watch
	popd
