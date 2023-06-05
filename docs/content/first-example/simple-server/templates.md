---
title: "Templates"
brief: "the implementation"
weight: 20
---

{{<lead>}}
Templates are the implementation for your generator.
They are parameterized files which are filled in
with data from the schema.
{{</lead>}}


Generators have several kinds of files that end up in the output

1. __Once Templates__ - used to generate a single file, like `main.go` or `index.js`
2. __Repeated Templates__ - generate a file for each element, like routes in this example
3. __Partial Templates__ - reusable template snippets which are available in all full templates
3. __Static Files__ - copied directly into the output, bypassing the template engine

Templates are based on Go `text/template` with extra helpers and conventions.
We will cover the basics in the first-example and they should be familiar to other text templating systems.
Read [template writing(/code-generation/template-writing/) to learn more about the details.

### Once Templates

These files are needed once for every server we generate.
Some have minimal templating and others loop over values, like `router.go`.

{{<codePane lang="text" title="templates/go.mod" file="code/first-example/simple-server/templates/go.mod" collapse="true">}}
{{<codePane lang="go" title="templates/server.go" file="code/first-example/simple-server/templates/server.go">}}

{{<codePane lang="go" title="templates/router.go" file="code/first-example/simple-server/templates/router.go">}}
{{<codePane lang="go" title="templates/middleware.go" file="code/first-example/simple-server/templates/middleware.go">}}

### Repeated and Partial Templates

We separate the handler into a template which uses the partial.
This is for demonstration purpose here and will be more useful
in the "full-example" section where the implementation is more complete.

{{<codePane lang="go" title="templates/route.go" file="code/first-example/simple-server/templates/route.go">}}
{{<codePane lang="go" title="partials/handler.go" file="code/first-example/simple-server/partials/handler.go">}}

### Rendered Output Files

Here we can see the result of code generation for a sample of the files.
We will actually generate these in the next section.
They are provided here so you can see the input / output pairs on a single page.

{{<codePane lang="go" title="output/middleware.go" file="code/first-example/simple-server/output/middleware.go">}}
{{<codePane lang="go" title="output/router.go" file="code/first-example/simple-server/output/router.go">}}
{{<codePane lang="go" title="output/routes/Hello.go" file="code/first-example/simple-server/output/routes/Hello.go">}}

