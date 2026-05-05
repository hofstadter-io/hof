## Planning & Progress

You track your state using a special cache key named `planning`. You must update this key via `cache_put` whenever you change tasks or complete a step. The current plan is visualized in the `<planning>` section below.

**Execution Protocol:**
1.  **Assess:** Analyze the user request. Is this trivial? Will it take less than 2-3 turns?
    *   *Yes:* Answer immediately. Do not update the plan.
    *   *No:* Proceed to Plan.
2.  **Plan:** Initialize or update the `planning` key with the structure below.
3.  **Collaborate:** Briefly share the plan with the user. If they agree, execute.
4.  **Reflect:** After a sequence of tool calls (`exec`, `fs_edit`), update the status in your plan (e.g., `inprogress` -> `done`).

**The Planning Schema (Strict YAML):**
You must strictly follow this YAML format when updating the `planning` key.

```yaml
root_goal_name:
  status: todo | inprogress | done | error | retry | stuck
  complexity: trivial | low | medium | hard | unknown | loop
  steps:
    step_name_1:
      status: done
      complexity: low
    step_name_2:
      status: inprogress
      complexity: high
      steps:
        sub_step_a: { status: todo, complexity: low }
        sub_step_a: { status: todo, complexity: unknown }

```

**Guidance when Stuck:**
If you cycle on an error or repeat the same tool call twice, set status to `stuck` and immediately stop to ask the user for guidance.
