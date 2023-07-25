---
title: "Schema"
brief: "for a REST server"
weight: 10
---

{{<lead>}}
__Hof Schemas__ are CUE values and serve as the contract between
you, the generator writer, and your users.
Your schemas capture the essence of the problem or application.
They provide the constraints for users developing with your generator
and ensure input to your generator is valid.
{{</lead>}}

### A Schema in Hof

A __Hof Schema__ is a CUE value.

{{<codePane title="schema/myschema.cue" file="code/first-example/simple-server/content/schema/minimal.html">}}


### A Minimal REST Schema

Let's start by sketching out the minimal definition for a server.
We put this in the `schema/` directory and thus CUE package.

{{<codePane title="schema/server.cue" file="code/first-example/simple-server/content/schema/rest.html">}}


### Routes and Handlers

Routes are a main focal point of REST servers.
When we generate the code, we will need handlers for each.
Here is the schema for a Route that will have a handler generated.

{{<codePane title="schema/server.cue" file="code/first-example/simple-server/content/schema/routes.html">}}


### Extra Features

These are features you may want to provide to your server users.
While the user only has to set a boolean or flag,
they can get advanced capabilities which are the consistent
for every generated server.

{{<codePane title="schema/server.cue" file="code/first-example/simple-server/content/schema/extra.html">}}


### Calculated Fields

These are fields and values you can infer from a user's input that they do not need to set.
They are often for making template writing easier by formatting or collecting values.
You can also transform or enrich user input using CUE here as well as in generators.
Later, we will use this technique to define the CRUD routes for a resource based on the user's data model.

{{<codePane title="schema/server.cue" file="code/first-example/simple-server/content/schema/calc.html">}}


### Language and Tool Fields

There are typically language or tools specifics which may need to be configured.
Often, you will want to make some of these accessible to your users.
You can even create generators which are multilingual or lets your
user select their preferred technologies.

{{<codePane title="schema/server.cue" file="code/first-example/simple-server/content/schema/lang.html">}}


### Final Schema

This is our final schema after combining all the parts above.

{{<codePane title="schema/server.cue" file="code/first-example/simple-server/schema/server.html">}}
