---
title: Docs

weight: 15
---


{{<lead>}}
This page provides an outline for developing and contributing
to `hof`'s documentation site.
{{</lead>}}

### Running

The documentation is under `docs/` and the following
assumes you are working from there.

Run `make first` once for development setup.

To run the website locally

- with drafts: `make dev`
- without drafts: `make stg`

Hugo's config is also defined as CUE (config.cue),
run `make config.yaml` to update.

GitHub Actions is used to build and deploy
for [next](https://next.hofstadter.io) and this site.


### Content

If you'd like to write non-trivial documentation content,
please open an issue or ask on Slack.
We are happy to help you outline a page or section.


Generally:

- There are many short-codes available in `layouts/shortcodes`.
- Each page should have the following fields set in the front matter:

```yaml
title: "..."
linkTitle: "..."
brief: "..."
description: "..."
keywords:
- kw1
- kw2

# pages are ordered by weight
weight: 1
```

(values are for demonstration)



### Example code

Most example code is under the `code/` directory
and referenced in the Markdown.

Use the following `shortcodes` for adding examples to Markdown files.

```go-html-template
{{</*codePane file="code/..." title="..." lang="..." */>}}

{{</* codePane2
	file1="code/..." title1="..." lang1="..."
	file2="code/..." title2="..." lang2="..."
*/>}}

{{</*codePane3
	file1="code/..." title1="..." lang1="..."
	file2="code/..." title2="..." lang2="..."
	faile3="code/..." title3="..." lang3="..."
*/>}}

{{</*codeInner title="..." lang="..."*/>}}
...
{{</*/codeInner*/>}}
```

CUE files require an extra step to generate highlighted code.

1. __Omit the lang arg for the shortcode when using CUE__
1. `make highlight` to turn CUE into HTML for embedding.
   These will be next to the CUE file after generating.



### Tests

We try to make our examples runnable and testable as much as possible.
You can see the Makefiles along side the examples under `code/`

Run `make verify` to test example code and ensure there are no output differences.

Run `make blc.dev` to check for broken links against localhost.



