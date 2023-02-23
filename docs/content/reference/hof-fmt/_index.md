---
title: "hof / fmt"
linkTitle: "hof / fmt"
weight: 40
---

{{<lead>}}
`hof/fmt` is a command which will
format any and all languages.
You can create your own formatters as well.
{{</lead>}}


`hof` needs a code formatter for the languages it generates.
It runs the pre-output through before applying diff and merging
with any custom code you added to output files.
This is simplifies the job of template authors,
but is also required to avoid unneccessary merge conflicts.

You will need Docker available to use this feature.
Hof will pull and run containers in the background.
You can disable this by setting an environment variable.

> `HOF_FORMAT_DISABLED=1`


{{<codePane file="code/cmd-help/fmt" title="$ hof help fmt" lang="text">}}
