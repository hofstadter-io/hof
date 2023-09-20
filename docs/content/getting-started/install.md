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

Git & (Docker,Podman,Nerdctl) should be available but are also optional.

## Installation

<br>

Recommended version: <b>{{<hof-rel-link>}}</b>

<br>

#### Binary downloads, rename the file to `hof` and place it in your PATH.

These are the same links for the curl.

{{<hof-dl-btns>}}

<br>

[All Releases](https://github.com/hofstadter-io/hof/releases)

[Container Images](https://github.com/orgs/hofstadter-io/packages?repo_name=hof)

<br>

#### Homebrew

<br>

{{<codeInner title="with Homebrew" lang="text">}}
brew install hofstadter-io/tap/hof
{{</codeInner>}}

<br>


#### Source 

<br>

{{<codeInner title="from source" lang="text">}}
git clone https://github.com/hofstadter-io/hof && cd hof
git checkout {{<hof-version>}}
make hof
{{</codeInner>}}

<br>







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



## Testing __hof__

Run `hof help` in your terminal.

{{<codePane file="code/cmd-help/hof" title="$ hof help" lang="text">}}



## Opting out of telemetry

We collect a few metrics to help us understand basic usage.

1. Which command is run, never the arguments or flags
2. OS & arch
3. the `hof` version

[link to code](https://github.com/hofstadter-io/hof/blob/_dev/cmd/hof/ga/ga.go#L153)

To opt out of telemetry, set `HOF_TELEMETRY_DISABLED=1`

Long term, we would like to move to the Go strategy, which
uses a probability based method for reducing metrics even further,
i.e to collect enough for statistical 
significance and much more focused to the questions at hand.

If you'd like to learn more about Go's telemetry proposal, see
https://github.com/golang/go/discussions/58409
