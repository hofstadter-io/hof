---
title: Resources

weight: 30
---

### Generator
- 
- additions to the generator, only needed to add files under ResourceFiles
- (notice how our additions have reduced from original, to datamodel, to client), put in a final section
- html without extension, hack so we can view contemporary url routes

{{< codePane title="gen/server.cue" file="code/first-example/web-client/content/gen/resource-web.html" >}}


### Initial HTML / JS

Since we broke down `index.html` into a number of partials,
we can reuse them here and end up with a template for resource html
which is very close to our index page.

{{< codePane title="templates/resource.html" lang="html" file="code/first-example/web-client/templates/resource.html" >}}

We will want to break this down into smaller snippets

- talk about routes & query, list vs item views

{{< codePane title="templates/resource.js" lang="js" file="code/first-example/web-client/content/templates/resource-prebreak.js" >}}


### Rendering Resources

- add templates for data rendering
- templates in templates (2 options) & custom code additions
    - using alt delims in output, string replace
		- using alt delims in template, gen config (maybe we put this later, mention that it can be done and provide a link)
- explain hack for using query (loc.search) for list vs item views

### Fracturing the template into partials

- break down JS template into partials
