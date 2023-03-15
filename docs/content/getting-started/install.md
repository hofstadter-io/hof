---
title: "Installation"
description: "Download and install the hof CLI tool"
brief: "Download and install the hof CLI tool"
weight: 3
---

{{<lead>}}
__hof__ is available for all major operating systems and architectures.
{{</lead>}}

git & docker should be available but are also optional

## Installation

<br>

Current version: <b>{{<hof-rel-link>}}</b>

<br>

{{<codeInner title="with Homebrew" lang="text">}}
brew install hof
{{</codeInner>}}

{{<codeInner title="as a binary" lang="text">}}
{{<hof-curl>}}
{{</codeInner>}}

{{<codeInner title="from source" lang="text">}}
go install github.com/hofstadter-io/hof/cmd/hof@latest
{{</codeInner>}}

<br>

#### Binary downloads, rename the file to `hof` and place it in your PATH.

these are the same links for the curl

{{<hof-dl-btns>}}

<br>

[All Releases](https://github.com/hofstadter-io/hof/releases)

[Container Images](https://hub.docker.com/r/hofstadter/hof/tags)




## Testing __hof__

Run `hof help` in your terminal.

{{<codePane file="code/cmd-help/hof" title="$ hof help" lang="text">}}



## Updating __hof__


The built-in update command can be used to check and install any version.

{{<codeInner lang="sh">}}
# Check for an update
hof update --check

# Update to the latest version
hof update

# Install a specific version
hof update --version vX.Y.Z
{{</codeInner>}}
