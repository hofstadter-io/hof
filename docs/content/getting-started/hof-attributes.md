---
title: "#hof & @attributes"

weight: 70
---

{{<lead>}}
The `#hof` definition is how generators, datamodels, and workflows are discovered.
Several attributes allow you to write shorthand for common settings.
{{</lead>}}


### The Attributes

By now, you have seen hof's attributes in the previous getting-started sections.
Hof turns these CUE attributes into `#hof` configuration.

- `@gen(<name>)` - the root of a generator
- `@datamodel(<name>)` - the root of a datamodel
- `@flow(<task>)` - the root of a workflow or a task type

Datamodels and workflows have a few more attributes
that can be specified under their root.
They are covered in the respective sections on each.

### Schema

The following is the schema for `#hof`.
Many generators, workflows, and datamodels
you import and use will fill this in for you.

You can embed the following schema in order to reference the content.
Hof will automatically inject it during the loading process to ensure correctness.

{{< codePane title="hof/schema.Hof" file="code/hof-schemas/hof.html" >}}
