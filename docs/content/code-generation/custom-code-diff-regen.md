---
title: "Custom Code, Diff, and Regen"

weight: 30
---


### Diff3 Mode for Customizing Output

`hof` has a diff engine so you can add custom code
to the outputs, regenerate, and not lose your work.

Add `--diff3` to enable this rendering mode.

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.cue schema.cue -T template.txt -O out/ --diff3
{{</codeInner>}}



