---
title: "Schema"
brief: "for a REST server"
weight: 10
---

{{<lead>}}
__Hof Schemas__ are CUE definitions and serve as the contract between
you, the generator writer, and your users.
Your schemas capture the essence of the problem,
set the constraints for users designing with your generator,
and are the input to your code generation templates.
{{</lead>}}


{{<codePane title="A minimal schema" file="code/first-example/simple-server/content/schema/minimal.html">}}


### A Minimal REST Schema

Let's start by sketching out the minimal definition for a server.
We put this in the `schema/` directory and thus CUE package.

{{<codePane title="schema/server.cue" file="code/first-example/simple-server/content/schema/rest.html">}}

Note how we use default values for "optional" fields in our generator UX,
rather than CUE optional syntax.
This avoids some nil-ness checks and conversions
and make template writing easier.


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
They are often for making template writing easier
by formatting or collecting values.

{{<codePane title="schema/server.cue" file="code/first-example/simple-server/content/schema/calc.html">}}


### Language Fields

There are typically language specifics which may need to be configured.
You will likely need to make some accessible to your users.

{{<codePane title="schema/server.cue" file="code/first-example/simple-server/content/schema/lang.html">}}


### Final Schema

<details>
<summary>Final Schema</summary>
{{<codePane title="schema/server.cue" file="code/first-example/simple-server/content/schema/final.html">}}
</details>
