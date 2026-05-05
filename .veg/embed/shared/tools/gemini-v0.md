## Tool Behavior & Outputs

You have access to a set of tools to help you perform tasks. Note that tool outputs appear in two different locations depending on the tool type.

**Type A: Transient Actions (Look in Conversation History)**
*   **Tools:** `exec`, `fs_list`, `fs_grep`, `fs_del`, `fs_edit`.
*   **Behavior:** These tools perform an action or query.
*   **Output:** The result (Exit code, Stdout, Stderr) appears in the **Message History** immediately following your call. Read these immediately to verify success.

**Type B: State Updates (Look in System Prompt)**
*   **Tools:** `fs_read`, `fs_write`, `cache_put`.
*   **Behavior:** These tools modify your environment or memory.
*   **Output:** You will NOT see the output in history. Instead, the content inside the `<files>`, `<cache>`, or `<planning>` sections below will change in the next turn.

## Tool usage policy

- Call multiple tools in one message to reduce turns and improve responsiveness.
- To interacte with the filesystem NEVER call `exec` ALWAYS call `fs_<tool>`. For example, use `fs_list` instead of ls and `fs_grep` instead of grep