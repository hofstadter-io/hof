---
title: "API Resources"
brief: "from Data Models"
weight: 30
---

With our types and data store in place,
we will now want to expose these in the API.

- Add a Resource definition to the schema
- Update our generator definition
- Add the template for API resources
- Wire resources into the router


### Schema Additions

We first add a schema for a resource
and a [CUE 'function'](https://cuetorials.com/patterns/functions/)
for converting models to resources.

{{< codePane title="gen/server.cue" file="code/first-example/data-layer/content/schema/resource.html" >}}

### Generator Changes

Add the following changes in their appropriate places into the existing generator definition.

{{< codePane title="gen/server.cue" file="code/first-example/data-layer/content/gen/resource.html" >}}


### Resource Template

The following creates CRUD handlers.
Note how we can reuse our route handler partial template
because we added these in the schema and mapping.

Create a new template called `resource.go`

{{< codePane title="templates/resource.go" file="code/first-example/data-layer/templates/resource.go" lang="go" >}}


### Other Templates

Some small changes to existing templates as well

{{<codePane title="templates/router.go" file="code/first-example/data-layer/content/templates/router.go" lang="go">}}

### Regenerate the Server

You can now run `hof gen ./example` and you should find a new `./output/resources` directory.

### Using the Resources

You could now rebuild and call our CRUD endpoints,
_except that we haven't yet implemented the handlers, which we will do in the next section_.

