  - Batch execution tool that runs multiple tool invocations in a single request
  - Tools are executed in parallel when possible, and otherwise serially
  - Takes a list of tool invocations (tool_name and input pairs)
  - Returns the collected results from all invocations
  - Use this tool when you need to run multiple independent tool operations at once -- it is awesome for speeding up your workflow, reducing both context usage and latency
  - Each tool will respect its own permissions and validation rules
  - The tool's outputs are NOT shown to the user; to answer the user's query, you MUST send a message with the results after the tool call completes, otherwise the user will not see the results
  
  Available tools:
  ${(
    await Promise.all(
      (await $c5()).map(
        async (Z) => `Tool: ${Z.name}
  Arguments: ${Rc5(Z.inputSchema)}
  Usage: ${await Z.prompt({ permissionMode: I })}`,
      ),
    )
  ).join(`
  ---`)}
  
  Example usage:
  {
    "invocations": [
      {
        "tool_name": "${c9.name}",
        "input": {
          "command": "git blame src/foo.ts"
        }
      },
      {
        "tool_name": "${rw.name}",
        "input": {
          "pattern": "**/*.ts"
        }
      },
      {
        "tool_name": "${uX.name}",
        "input": {
          "pattern": "function",
          "include": "*.ts"
        }
      }
    ]
  }
