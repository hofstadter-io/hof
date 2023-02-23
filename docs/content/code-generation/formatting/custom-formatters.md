---
title: Custom Formatters
weight: 30
---

{{<lead>}}
Most of Hof's code formatters are implemented as a simple server in a container.
This makes it easy to build and use your own.
{{</lead>}}


To learn how to build a formatter, see [hof/formatters](https://github.com/hofstadter-io/hof/tree/_dev/formatters) on GitHub.

To use a custom formatter, configure it in a generator using the [Formatting](https://github.com/hofstadter-io/hof/blob/_dev/schema/gen/generator.cue#L32) section.

For file specific overrides, use the [Formatting](https://github.com/hofstadter-io/hof/blob/_dev/schema/gen/file.cue#L32) value.
