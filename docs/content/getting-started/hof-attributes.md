---
title: "#hof & @attributes"
brief: "Special annotations for hof."

weight: 70
---

{{<lead>}}
The `#hof` definition is how generators, datamodels, and workflows are discovered and configured.
Several attributes allow you to write shorthand for common settings.
{{</lead>}}


### The Attributes

By now, you have seen hof's attributes in the previous getting-started sections.
Hof turns these CUE attributes into `#hof` configuration.

- `@gen(<name>)` - the root of a generator
- `@datamodel(<name>)` - the root of a datamodel
- `@flow(<name>)` - the root of a workflow

Datamodels and workflows have attributes that can be used under their root
and are covered in their respective sections.

Other attributes:

- `@userfiles(<glob>)` - load file contents into a struct as strings, the key is the filepath.
  This works with CUE based commands.

### Schema

The following is the schema for `#hof`.
Many generators, workflows, and datamodels
you import and use will fill this in for you.

You can embed the following schema in order to reference the content.
Hof will automatically inject it during the loading process to ensure correctness.

{{< codePane title="hof/schema.Hof" file="code/hof-schemas/hof.html" >}}
