---
title: "Installation"
description: "Download and install the hof cli tool"
brief: "Download and install the hof cli tool"
weight: 3
---

{{<lead>}}
__hof__ is available for all major operation systems and architectures.
{{</lead>}}

git & docker should be available, but are also optional

## Installation

<br>

Current version: <b>{{<hof-rel-link>}}</b>

<br>

{{<codeInner title="installation commands" lang="text">}}
// with homebrew
brew install hof


// as a binary
{{<hof-curl>}}

// or from source
go install github.com/hofstadter-io/hof/cmd/hof@latest
{{</codeInner>}}

<br>

#### Binary downloads, rename the file to `hof` and place it in your PATH.

these are the same links for curl

{{<hof-dl-btns>}}

<br>

[All Releases](https://github.com/hofstadter-io/hof/releases)

[Container Images](https://hub.docker.com/r/hofstadter/hof/tags)




## Testing __hof__

Run `hof help` in your terminal.

{{<codePane file="code/cmd-help/hof" title="$ hof help" lang="text">}}



## Updating __hof__


The builtin update command can be used to check and install any version.

{{<codeInner lang="sh">}}
# Check for an update
hof update --check

# Update to the latest version
hof update

# Install a specific version
hof update --version vX.Y.Z
{{</codeInner>}}


