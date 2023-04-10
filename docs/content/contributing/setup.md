---
title: Setup

weight: 5
---


{{<lead>}}
This page will help you get your development environment setup
and commands for working on `hof` and the documentation.
{{</lead>}}


### Tools

For `hof`: CUE, Docker, Make

Requires CUE >= 0.5.0
[using one of the methods here](https://cuelang.org/docs/install/)


For `docs`: Hugo, npm

Requires Hugo >= 0.111, the extended version.
You can [download it from GitHub](https://github.com/gohugoio/hugo/releases)

We also have a GitHub Codespace with the tools installed.
You can launch this from the repository.


### GitHub Actions

Our CI runs in GitHub Actions.
We define the actions as CUE and then
generate the yaml into `.github/workflows`

- The CUE is in `ci/gha`
- Run `make workflow` to generate the yaml


### ENV and debugging settings


There are a few flags, environment variables, and debug settings to be aware of.

- `--verbose N` is used to increase printed messages when running hof
- Depending on the tests you wish to run, you may need various ENV VARs setup. (notably `hof mod`)
- Several core packages have a `debug bool` variable that can be set to true


