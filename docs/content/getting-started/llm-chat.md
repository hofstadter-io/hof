---
title: LLM Chat
description: "Combining LLMs, Bard, ChatGPT, and Hof."
brief: "Combining LLMs, ChatGPT, and hof."

keywords:
- LLM
- Bard 
- ChatGPT
- code gen

weight: 50
---

{{<lead>}}
Large Language Models (LLM) are an inflection point in computing.
The represent significant advancements and automation for tasks.
Generating code is among them and there are many interesting topics
at the intersection of LLMs and Hof.
{{</lead>}}

{{<alert style="info">}}
We are only at the beginnings of our merging Hof with LLMs.
This page descibes the current state and where we are headed.
{{</alert>}}


## hof chat

The `hof chat` command is and early preview for interacting with hof using natural language prompts.
You can already use this to:

1. Talk with ChatGPT from the command line or vim
2. Talk with Hof data models (full demo coming soon :)

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

