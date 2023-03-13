---
title: "Creators"
description: "Application blueprints and bootstrapping from in any git repo"
brief: "Application blueprints and bootstrapping from in any git repo"

weight: 30
---


{{<lead>}}
Hof Creators are code generator modules
intended for one-time bootstrapping from a few inputs.
Use them to give your users one-liners for getting started
or to provide reusable application blueprints.
{{</lead>}}

{{<editorNote editor=sara color=red >}}
Hey Tony!

blah blah blah
{{</editorNote>}}

## Running `hof create`

`hof create` is run by users who provide the url to a git repository.

{{<codeInner title="> terminal">}}
$ hof create github.com/hofstadter-io/hofmod-cli
{{</codeInner>}}

The user is presented with an interactive prompt you design.
Their answers provide the input to your generator,
which should output the files they need to get started using your project. 
The process is flexible enough for any git repository to provide the creator experience,
so what you generate does not have to be `hof` related.

<div id="create-cast" class="asciinema"></div>

<br>
<br>



## Adding a Creator to your repository

To add a creator to a git repository, we have to set up
a CUE module and hof generator, and then push a tag.

1. initialize a CUE repository in your project (`hof mod init cue <repo-url>`)
1. add a `creator.cue` at the root of your repository to hold the generator
1. fill in the generator, templates, and prompt
1. test the creator locally
1. tag and push

There are a number of options for how to
organize a creator, provide the prompt and inputs,
and conditionally ask questions or output files.
More information is provided below.

## Example Creator

The following is the creator from the `hofmod-cli` generator
in the video above.

{{<codePane file="code/hof-create/hofmod-cli/creator.html" title="hofmod-cli/creator.cue">}}

## Learn More

As we were saying, a `creator` is a hof generator with a `Create` field,
run as a one-time code generation.
The main new part is the `prompt` where you can ask the user for inputs.
These then form the input value to code generation,
and at this point, the creator is like any other generator

Users then run `hof create <repo>@<tag>` to bootstrap
full applications or features to an existing application.
Since creators are just generators, you can create the files
needed for CI/CD, deployment, security, or other system.
You can output CUE files and many `hof generators`
will also provide creators to bootstrap the
initial files needed to use their module.

To learn more about writing creators, prompts, and how to give your users
one-line bootstrapping or application blueprints, see the
[hof create section](/code-generation/creators/)


## Creator vs Generator Module

Under the hood, a Creator is a one-time Generator Module
with an optional prompt.
To decide when to use one versus the other...

Use a Creator when:

- You expect users to run the generation process just once
- Your users need to quickly bootstrap an application based on a few inputs
- You want a `hof create github.com/user/repo` one-liner to give your users

Use a Module when:

- You expect your users to run code generation iteratively, as part of their application development
- You want your users to add or change features and then regenerate their code
- If you want to give them the ability to change or add features from your app skeleton, then hof gen is the way to go.
  Full generators were designed for an iterative, ongoing experience during software development.


