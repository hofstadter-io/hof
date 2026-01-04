DAG_VARS=_EXPERIMENTAL_DAGGER_RUNNER_HOST=docker-container://veg-dagger-engine
query:
	echo '{engine{localCache{entrySet}}}' | sh -c '$(DAG_VARS) dagger query > dagger-cache.json'

inspect:
	cue eval dagger-cache.cue dagger-cache.json -e e2 -e total