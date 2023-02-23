---
title: "Template Writing"
weight: 20
---

## Template System

Hof uses Go's `text/template` package (https://pkg.go.dev/text/template).
All partials are registered with all templates for use with `{{ template "path/to/partial.ext" }}`.

#### Sections:

{{< childpages >}}

## Notes

- how to write templates to avoid merge conflicts with diff3 algos
  - don't have users write code at the edge of a range in template
	- add a comment like `// edit below` or `// edit between`
	- need stable code on both sides of where users will write
- giving your users a function or stub to fill in, keep separate
  - isolates custom code
	- gen code will call this
	- gen code can wrap / transform around this so users don't have to update
	- this will be where we put the lenses type capabilities
- templates generally
- loaded partials
- inline partials
- page on how to use go's templates

## Debugging

The best way is to use comments and the yaml helper to inspect values.
You can reduce, and then build back up, any conditional or loop logic.
Often, a separate, development only, file will be used for this.

{{<codeInner lang="go">}}
/*
{{ yaml .MyValue }}
*/
{{</codeInner>}}


