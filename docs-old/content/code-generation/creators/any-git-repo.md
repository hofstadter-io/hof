---
title: "Hof Create for any git repository"
linkTitle: "any git repo"
weight: 100
---

{{<lead>}}
This section will show you how to add a generator
to any repository, hof generator or not.
The process is the same, the only difference
is what files you bootstrap for your users.
{{</lead>}}

1. Write a small creator
1. Push a git tag
1. `hof create github.com/user/repo@<tag>`
1. Prompt the user for a few inputs
1. Generate some files for your user


### developing

Copy the example from the earlier section
or create your own into a `name.cue` file.

