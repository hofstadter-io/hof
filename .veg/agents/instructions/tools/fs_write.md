adds or overwrites an the file at $path with the $content, the changes will be available in your working key/value cache and context next turn

IMPORTANT: Do not call `fs_read` to verify changes were applied to the the contents. 
Changes will be automatically updated to both the file system and your cache, making them available immediatly in your system prompt the next turn.
Trust the process and that the file's contents in your system prompt are accurate and up-to-date.