---
title: "Installation"
description: "Download and install the hof CLI tool."
brief: "Download and install the hof CLI tool."
keywords: "homebrew, cli, get started, code generation tool, installation guide, DAG flow, run CUE pipelines, CUE dependency management, go mods"
weight: 3
---

{{<lead>}}
__hof__ is available for all major operating systems and architectures.
{{</lead>}}

Git & Docker should be available but are also optional.

## Installation

<br>

Current version: <b>{{<hof-rel-link>}}</b>

<br>

{{<codeInner title="with Homebrew" lang="text">}}
brew install hofstadter-io/tap/hof
{{</codeInner>}}

{{<codeInner title="as a binary" lang="text">}}
{{<hof-curl>}}
{{</codeInner>}}

{{<codeInner title="from source" lang="text">}}
go install github.com/hofstadter-io/hof/cmd/hof@latest
{{</codeInner>}}

<br>

#### Binary downloads, rename the file to `hof` and place it in your PATH.

These are the same links for the curl.

{{<hof-dl-btns>}}

<br>

[All Releases](https://github.com/hofstadter-io/hof/releases)

[Container Images](https://hub.docker.com/r/hofstadter/hof/tags)




## Testing __hof__

Run `hof help` in your terminal.

{{<codePane file="code/cmd-help/hof" title="$ hof help" lang="text">}}



## Updating __hof__


You can use the built-in update command to check and install any version.

{{<codeInner lang="sh">}}
# Check for an update
hof update --check

# Update to the latest version
hof update

# Install a specific version
hof update --version vX.Y.Z
{{</codeInner>}}
