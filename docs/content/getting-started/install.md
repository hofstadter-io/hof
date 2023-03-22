---
title: "Installation"
description: "Download and install the hof CLI tool."
brief: "Download and install the hof CLI tool."
keywords:
  - homebrew 
  - downloads
  - update
weight: 3
---

{{<lead>}}
__hof__ is available for all major operating systems and architectures.
{{</lead>}}

Git & Docker should be available but are also optional.

## Installation

<br>

Current version: <b>{{<hof-rel-link>}}</b>

{{< alert style="warning" >}}
We recommend the latest beta over v0.6.7 as
the core features have seen many improvements.
You can expect it to be more stable and correct.
{{</alert>}}

<br>

{{<codeInner title="as a binary" lang="text">}}
{{<hof-curl>}}
{{</codeInner>}}

{{<codeInner title="from source" lang="text">}}
go install github.com/hofstadter-io/hof/cmd/hof@latest
{{</codeInner>}}

{{<codeInner title="with Homebrew" lang="text">}}
brew install hofstadter-io/tap/hof
{{</codeInner>}}


#### Binary downloads, rename the file to `hof` and place it in your PATH.

These are the same links for the curl.

{{<hof-dl-btns>}}

<br>

[All Releases](https://github.com/hofstadter-io/hof/releases)

[Container Images](https://github.com/orgs/hofstadter-io/packages?repo_name=hof)




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

with Homebrew

{{<codeInner lang="sh">}}
# Check for an update in brew
brew outdated hofstadter-io/tap/hof

# Update to the latest version in brew
brew upgrade hofstadter-io/tap/hof

# To get more info regarding hof package
brew info hofstadter-io/tap/hof
{{</codeInner>}}
