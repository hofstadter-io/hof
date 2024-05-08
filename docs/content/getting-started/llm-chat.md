---
title: LLM Chat
description: "Combining LLMs and Hof."
brief: "Combining LLMs and hof."

keywords:
- LLM
- Gemini
- ChatGPT
- code gen

weight: 50
---

{{<lead>}}
Large Language Models (LLM) are an inflection point in computing.
`hof chat` is an experiment in combining LLMs with hof's features
to find where and how they can be used together for the best effect.
{{</lead>}}


## hof chat

The `hof chat` command is and early preview for interacting with hof using natural language prompts.
You can already use this to:

1. Talk with ChatGPT from the command line or vim
1. Talk with Hof data models (full demo coming soon :)

_note, you can also call any LLM apis via hof/flow to build complex workflows_

<br>

{{<codePane file="code/cmd-help/chat" title="$ hof help chat" lang="text">}}

## where we are going

We see Hof + LLM as better than either on their own.

__LLMs provide for natural language interfaces to all things Hof__

We are building a future where LLM powered Hof is your coding assistant,
allowing you to use the best interface (LLM, IDE, low-code) for the task at hand.

__Hof simplifies code gen with LLMs__

Hof's deterministic code gen means that the LLMs only have to generate the
data models and extra configuration needed for generators. This has many benefits.

- The task for the LLM is much easier and they can do a much better job.
- The code generation is backed by human written code, so no hallucinations.
- The same benefits for generating code at scale with Hof.

Other places we see LLMs helping Hof

- importing existing code to CUE & Hof
- automatically transforming existing code to hof generators
- filling in the details and gaps in generated code
- in our premium user interfaces for low-code
  (these are more the multi-modal models, which come after LLMs, think Google Gemini)

