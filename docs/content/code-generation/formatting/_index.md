---
title: Formatting
weight: 25
---

{{<lead>}}
Consistent formatting makes code more readable and understandable.
Hof automatically formats many languages with widely used defaults
while also providing a means to customize or add new languages.
{{</lead>}}

When you generate code with `hof` it will automatically format
the output for known extensions and languages.
This typically happens by running a container in the background,
though some languages are handled directly in `gocode`. (See below)

`hof` will start a container the first time it is needed for you.
You can fetch, run, and stop formatters manually with the [hof fmt](/code-generation/formatting/format-command/).

We have built several [hof/formatters](https://github.com/hofstadter-io/hof/tree/_dev/formatters)
which use the most common tool for a language and use the most commonly used rules.
You can build your own and configure generators to use them too,
see the [custom formatters](/code-generation/formatting/custom-formatters/) section

```
hof fmt
```

{{<lead>}}
`hof/fmt` is a command which will
format any and all languages.
You can create your own formatters as well.
{{</lead>}}


`hof` needs a code formatter for the languages it generates.
It runs the pre-output through before applying diff and merging
with any custom code you added to output files.
This is simplifies the job of template authors,
but is also required to avoid unnecessary merge conflicts.

You will need Docker available to use this feature.
Hof will pull and run containers in the background.
You can disable this by setting an environment variable.

> `HOF_FMT_DISABLED=1`


{{<codePane file="code/cmd-help/fmt" title="$ hof help fmt" lang="text">}}

## Supported Languages

If a language you need is missing, [open an issue](https://github.com/hofstadter-io/hof/issues/new?title=hof%2Ffmt%3A%20add%20support%20for%20%3Clanguage%3E).
If you can, please provide a link to a common tool used for that language.

<br>


| Language | Tool |
|----------|:----:|
| csharp   | csharpier |
| css |     prettier/css |
| cue      | gocode |
| go       | gocode |
| graphql | prettier/graphql |
| groovy |  prettier/groovy |
| html |    prettier/html |
| java |    prettier/java |
| js |      prettier/babel |
| json     | gocode |
| jsx |     prettier/babel |
| less |    prettier/less |
| md |      prettier/markdown |
| php |     prettier/php |
| python   | black |
| rb |      prettier/ruby |
| rs |      prettier/rust |
| scss |    prettier/scss |
| toml     | gocode |
| ts |      prettier/typescript |
| tsx |     prettier/typescript |
| vue |     prettier/vue |
| xml      | gocode |
| yaml     | gocode |


