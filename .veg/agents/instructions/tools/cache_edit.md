This is a tool for editing files, snippets, or values in the cache.

Before using this tool:

1. Make sure you understand the file or values's contents and context.
2. If a file is not in the cache, use cache_file to load it.
3. If you want to make a new file or replace all contents, use cache_write instead
4. Your cache edits will be available your next turn

When making multiple file edits in a row to the same file, or similar/related edits across multiple files,
always do so by sending a single message/turn with multiple calls to this edit tool. 

When making edits:
    - Ensure the edit results in idiomatic, correct code
    - Do not leave the code in a broken state
    - Always use absolute file paths (starting with /)

To make a cache edit, provide the following:

1. file_path: The absolute path to the file to modify (must be absolute, not relative)
2. old_string: The text to replace (must match the file contents exactly, including all whitespace and indentation)
3. new_string: The edited text to replace the old_string
4. expected_replacements: The number of replacements you expect to make. Defaults to 1 if not specified.

By default, the tool will replace ONE occurrence of old_string with new_string in the specified file. If you want to replace multiple occurrences, provide the expected_replacements parameter with the exact number of occurrences you expect.

CRITICAL REQUIREMENTS FOR USING THIS TOOL:

1. UNIQUENESS (when expected_replacements is not specified): The old_string MUST uniquely identify the specific instance you want to change. This means:
    - Include AT LEAST 3-5 lines of context BEFORE the change point
    - Include AT LEAST 3-5 lines of context AFTER the change point
    - Include all whitespace, indentation, and surrounding code exactly as it appears in the file

2. EXPECTED MATCHES: If you want to replace multiple instances:
    - Use the expected_replacements parameter with the exact number of occurrences you expect to replace
    - This will replace ALL occurrences of the old_string with the new_string
    - If the actual number of matches doesn't equal expected_replacements, the edit will fail
    - This is a safety feature to prevent unintended replacements

3. VERIFICATION: Before using this tool:
    - Check how many instances of the target text exist in the file
    - If multiple instances exist, either:
      a) Gather enough context to uniquely identify each one and make separate calls, OR
      b) Use expected_replacements parameter with the exact count of instances you expect to replace

WARNING: If you do not follow these requirements:
    - The tool will fail if old_string matches multiple locations and expected_replacements isn't specified
    - The tool will fail if the number of matches doesn't equal expected_replacements when it's specified
    - The tool will fail if old_string doesn't match exactly (including whitespace)
    - You may change unintended instances if you don't verify the match count

Remember: When making multiple file edits in a row to the same file, or similar/related edits across multiple files,
always do so by sending a single message/turn with multiple calls to this edit tool. 

