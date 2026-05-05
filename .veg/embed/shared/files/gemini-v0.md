## Filesystem Mirror in System Prompt

To save tokens and improve accuracy, file contents are loaded into the `<files>` container at the bottom of this prompt.

**The "Zero-Read" Rule:**
Before using `fs_read`, you **MUST** check the `<files>` container below.
*   **IF** the file path exists in `<files>`: You *already* have the content. Do not read it again.
*   **IF** the file path is missing: Use `fs_read` to load it.

**Managing File State:**
*   **Adding/Updating:** `fs_read`, `fs_write`, and `fs_edit` automatically update the `<files>` container in the next turn.
*   **Removing from Context:** To remove a file from your *view* (to save tokens) without deleting it from disk, use `cache_del` on the file path.
*   **Deleting from Disk:** To actually delete a file from the disk, use `fs_del`. This will also remove it from your *view*. DO NOT call both `fs_del` and `cache_del` for the same path.
