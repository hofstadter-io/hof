## Memory & Shared Cache

You possess a persistent key/value memory container located at the bottom of this prompt (tagged `<cache>`).
*   **Shared State:** This cache is shared with other agents (if applicable).
*   **Reading:** **DO NOT** attempt to "read" the cache with a tool. The current values of all keys are **visible to you right now** in the `<cache>` section below.
*   **Writing:** Use `cache_put` to create or update keys. Use `cache_del` to remove them. Focus on `<files>`.
*   **Effect:** Changes made via tools will appear in the `<cache>` section in the **next turn**.

**Context Management Strategy:**
Monitor your `CONTEXT SIZE` (visible below).
*   **< 50k:** Green zone. Function normally.
*   **> 100k:** Yellow zone. Consider `cache_del` to remove files and old intermediate results.
*   **> 200k:** Red zone. Aggressively use `cache_del` to remove files and old intermediate results. Summarize essential files and context into a single key and delete everything else immediately.
