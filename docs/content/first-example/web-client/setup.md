---
title: Setup
brief: "index page and wiring up the initial web client"

weight: 10
---

In this section, we will start by adding an `index.html`
and wiring our server to server static content.
We want a tracer bullet to for our web client.
This will get the process of generating and serving
our web UI in place. We will then add in more complexity
for resources and js to get, set, and display data.

(explain tracer bullet more)
(add link to webpage about it)

### Generator

We are going to add an `index.html` file, so first
we make sure it will be generated from a template.
We update our server generator to add `index.html`
to the `OnceFiles`.

{{< codePane title="gen/server.cue" file="code/first-example/web-client/content/gen/index-html.html" >}}

### Template

Create an `index.html`. This displays a minimal welcome message.

{{< codePane title="templates/index.html" lang="html" file="code/first-example/web-client/content/templates/index-init.html" >}}


### Serving

In order to serve our HTML, we need
to adjust our API routes.
First we need to move our existing routes to `/api`,
to make room for our static content route
Second, we need to server our static content from `/`
so index.html behaves correctly in the browser.

{{< codePane title="templates/router.go" lang="go" file="code/first-example/web-client/content/templates/router.go" >}}

### Regen, Rebuild, Rerun

(the three Re's)


### Breaking down index into partials

- break down now, prep for next section and reuse
- add navbar while we are at it
