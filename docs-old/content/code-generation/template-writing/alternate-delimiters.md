---
title: Alternate Delimiters
weight: 30
---

{{<lead>}}
Hof uses Go's text/template package which defaults to `{{` & `}}` for delimiters.
This can conflict with other systems, so hof enables you to change them.
{{</lead>}}


### Per-file Template Delimiters

You can configure delimiters per-file
by setting the `TemplateDelims` field
on a `gen.#File`.

{{< codePane title="per-file template delimiters" file="code/code-generation/template-writing/delimiters/per-file.html" >}}


### Per-glob Template Delimiters

You can set delimiters when loading
`Templates` and `Partials` when
you have a directory of files
which need alternative delimiters,
or you prefer a different syntax.
Set the `Delims` field in the configuration.

{{< codePane title="per-glob template delimiters" file="code/code-generation/template-writing/delimiters/per-glob.html" >}}
