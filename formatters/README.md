# fmt

a collection of code formatting tools

### formatting

- go, cue, data formats built in
- container based api wrappers for other languages
- run in background, avoid startup costs (dagger does this)
- build a container, config for control, where to find / how to run tools

links:

- https://github.com/rishirdua/awesome-code-formatters#general-purpose
- https://github.com/Unibeautify/docker-beautifiers
- https://github.com/NiklasPor/prettier-plugin-go-template
- https://prettier.io/docs/en/api.html
- https://prettier.io/docs/en/plugins.html


### Config a custom formatter

(tbd) enable formatter config through gens / files

### Writing a custom formatter

Create a folder under `./tools` with the name of your formatter

Make an API which, on the root route `/`

- input: `{ config: _, source: string }`
- output:
  - 200 formatted text
	- 400 formatted error

#### Share publicly

docker push your image, you (tbd) can now configure your gens to use it

#### Commit to hof

hook into `./lib/fmt/*.go`, see existing languages for examples

open a PR, get questions answered on slack


### Mutliarch container support

We need to use `docker buildx` or `dagger`, and hardware / vm setup 
which also supports qemu being setup right.
Maybe we just build them in github actions going forward since they
have to be pushed and then pulled anyway...

Searches related to building in GHA

- https://www.google.com/search?q=docker+multi+arch+builds+github+actions
- https://blog.thesparktree.com/docker-multi-arch-github-actions

Search results when building locally, never got it working

- https://stackoverflow.com/questions/60080264/docker-cannot-build-multi-platform-images-with-docker-buildx
- https://github.com/docker/buildx/issues/495

